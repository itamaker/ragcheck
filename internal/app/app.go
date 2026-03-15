package app

import (
	"fmt"
	"os"
)

func Run(args []string) int {
	if len(args) == 0 {
		return runTUI()
	}

	switch args[0] {
	case "score":
		return runScore(args[1:])
	case "judge":
		return runJudge(args[1:])
	case "tui", "interactive":
		return runTUI()
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand %q\n\n", args[0])
		usage()
		return 2
	}
}

func usage() {
	fmt.Println("ragcheck scores retrieval runs and judges RAG answers.")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  ragcheck                   # launch Bubble Tea TUI")
	fmt.Println("  ragcheck score -qrels examples/qrels.json -run examples/run.json -k 3")
	fmt.Println("  ragcheck judge -input examples/judge.json")
}
