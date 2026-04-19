package main

import (
	"log"

	"github.com/haukened/mirrorselect/internal/app"
)

func main() {
	cmd := app.NewRootCmd()
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
