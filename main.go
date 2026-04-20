package main

import (
	"log"
	"os"

	"github.com/haukened/mirrorselect/internal/app"
)

func main() {
	cmd := app.NewRootCmd()
	if err := cmd.Execute(); err != nil {
		if exitCode, ok := app.SudoReexecExitCode(err); ok {
			os.Exit(exitCode)
		}
		log.Fatal(err)
	}
}
