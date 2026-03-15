package main

import (
	"os"

	"github.com/itamaker/ragcheck/internal/app"
)

func main() {
	os.Exit(app.Run(os.Args[1:]))
}
