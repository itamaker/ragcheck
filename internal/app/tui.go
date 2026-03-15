package app

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type tuiField struct {
	key         string
	label       string
	placeholder string
	optional    bool
}

type tuiAction struct {
	name        string
	description string
	fields      []tuiField
	run         func(values map[string]string) (string, error)
}

type tuiState int

const (
	tuiMenu tuiState = iota
	tuiForm
	tuiResult
)

type tuiModel struct {
	title       string
	description string
	actions     []tuiAction
	state       tuiState
	menuInput   textinput.Model
	fieldInput  textinput.Model
	selected    int
	fieldIndex  int
	values      map[string]string
	output      string
	err         string
}

func runTUI() int {
	model := newTUIModel("ragcheck", "Interactive retrieval scoring and RAG judging", buildTUIActions())
	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func newTUIModel(title string, description string, actions []tuiAction) tuiModel {
	menu := textinput.New()
	menu.Placeholder = "Enter action number"
	menu.Focus()
	menu.CharLimit = 3
	menu.Width = 24

	field := textinput.New()
	field.CharLimit = 256
	field.Width = 64

	return tuiModel{
		title:       title,
		description: description,
		actions:     actions,
		state:       tuiMenu,
		menuInput:   menu,
		fieldInput:  field,
		values:      map[string]string{},
	}
}

func (m tuiModel) Init() tea.Cmd { return textinput.Blink }

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case tuiMenu:
			return m.updateMenu(msg)
		case tuiForm:
			return m.updateForm(msg)
		case tuiResult:
			return m.updateResult(msg)
		}
	}
	return m, nil
}

func (m tuiModel) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "enter":
		index := parseMenuSelection(m.menuInput.Value())
		if index < 0 || index >= len(m.actions) {
			m.err = "Invalid selection"
			return m, nil
		}
		m.selected = index
		m.state = tuiForm
		m.fieldIndex = 0
		m.values = map[string]string{}
		m.output = ""
		m.err = ""
		m.fieldInput = textinput.New()
		m.fieldInput.CharLimit = 256
		m.fieldInput.Width = 64
		m.fieldInput.Placeholder = m.actions[index].fields[0].placeholder
		m.fieldInput.Focus()
		return m, textinput.Blink
	}

	var cmd tea.Cmd
	m.menuInput, cmd = m.menuInput.Update(msg)
	return m, cmd
}

func (m tuiModel) updateForm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	action := m.actions[m.selected]
	field := action.fields[m.fieldIndex]

	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "esc":
		m.state = tuiMenu
		m.menuInput.SetValue("")
		m.menuInput.Focus()
		m.err = ""
		return m, textinput.Blink
	case "enter":
		value := strings.TrimSpace(m.fieldInput.Value())
		if value == "" && !field.optional {
			m.err = field.label + " is required"
			return m, nil
		}
		m.values[field.key] = value
		m.err = ""
		if m.fieldIndex == len(action.fields)-1 {
			output, err := action.run(m.values)
			m.output = strings.TrimSpace(output)
			if err != nil {
				m.err = err.Error()
			}
			m.state = tuiResult
			return m, nil
		}
		m.fieldIndex++
		next := action.fields[m.fieldIndex]
		m.fieldInput.SetValue("")
		m.fieldInput.Placeholder = next.placeholder
		return m, textinput.Blink
	}

	var cmd tea.Cmd
	m.fieldInput, cmd = m.fieldInput.Update(msg)
	return m, cmd
}

func (m tuiModel) updateResult(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "b":
		m.state = tuiMenu
		m.menuInput.SetValue("")
		m.menuInput.Focus()
		m.output = ""
		m.err = ""
		return m, textinput.Blink
	}
	return m, nil
}

func (m tuiModel) View() string {
	var b strings.Builder
	b.WriteString(m.title + "\n")
	b.WriteString(m.description + "\n\n")

	switch m.state {
	case tuiMenu:
		b.WriteString("Actions:\n")
		for i, action := range m.actions {
			b.WriteString(fmt.Sprintf("  %d. %s\n     %s\n", i+1, action.name, action.description))
		}
		b.WriteString("\nSelect an action: " + m.menuInput.View() + "\n")
		b.WriteString("Enter to continue, q to quit.\n")
	case tuiForm:
		action := m.actions[m.selected]
		field := action.fields[m.fieldIndex]
		b.WriteString("Action: " + action.name + "\n")
		b.WriteString(action.description + "\n\n")
		if len(m.values) > 0 {
			b.WriteString("Collected inputs:\n")
			for _, f := range action.fields[:m.fieldIndex] {
				b.WriteString(fmt.Sprintf("  - %s: %s\n", f.label, m.values[f.key]))
			}
			b.WriteString("\n")
		}
		label := field.label
		if field.optional {
			label += " (optional)"
		}
		b.WriteString(label + ": " + m.fieldInput.View() + "\n")
		b.WriteString("Enter to continue, esc to go back, q to quit.\n")
	case tuiResult:
		if m.err != "" {
			b.WriteString("Error:\n" + m.err + "\n\n")
		}
		if m.output != "" {
			b.WriteString("Output:\n" + m.output + "\n\n")
		}
		b.WriteString("Press b to go back or q to quit.\n")
	}

	if m.err != "" && m.state != tuiResult {
		b.WriteString("\nError: " + m.err + "\n")
	}
	return b.String()
}

func buildTUIActions() []tuiAction {
	return []tuiAction{
		{
			name:        "score",
			description: "Compute retrieval metrics from qrels and run files",
			fields: []tuiField{
				{key: "qrels", label: "Qrels Path", placeholder: "examples/qrels.json"},
				{key: "run", label: "Run Path", placeholder: "examples/run.json"},
				{key: "k", label: "Top-K", placeholder: "5", optional: true},
				{key: "json", label: "JSON Output", placeholder: "false", optional: true},
			},
			run: func(values map[string]string) (string, error) {
				args := []string{"-qrels", values["qrels"], "-run", values["run"]}
				if values["k"] != "" {
					args = append(args, "-k", values["k"])
				}
				if truthy(values["json"]) {
					args = append(args, "-json")
				}
				return captureRun(runScore, args)
			},
		},
		{
			name:        "judge",
			description: "Judge answer grounding and relevance from a record file",
			fields: []tuiField{
				{key: "input", label: "Judge Input Path", placeholder: "examples/judge.json"},
				{key: "json", label: "JSON Output", placeholder: "false", optional: true},
			},
			run: func(values map[string]string) (string, error) {
				args := []string{"-input", values["input"]}
				if truthy(values["json"]) {
					args = append(args, "-json")
				}
				return captureRun(runJudge, args)
			},
		},
	}
}

func captureRun(fn func([]string) int, args []string) (string, error) {
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	reader, writer, err := os.Pipe()
	if err != nil {
		return "", err
	}
	defer reader.Close()

	os.Stdout = writer
	os.Stderr = writer
	code := fn(args)
	_ = writer.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	body, readErr := io.ReadAll(reader)
	if readErr != nil {
		return "", readErr
	}
	output := string(body)
	if code != 0 {
		message := strings.TrimSpace(output)
		if message == "" {
			message = fmt.Sprintf("command exited with code %d", code)
		}
		return output, errors.New(message)
	}
	return output, nil
}

func parseMenuSelection(value string) int {
	choice, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || choice <= 0 {
		return -1
	}
	return choice - 1
}

func truthy(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "t", "yes", "y":
		return true
	default:
		return false
	}
}
