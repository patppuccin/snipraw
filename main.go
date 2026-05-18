package main

import (
	"os"

	"github.com/patppuccin/snipraw/src/cmd"
	"github.com/patppuccin/snipraw/src/console"
)

func main() {
	if err := cmd.SRCmd.Execute(); err != nil {
		console.Error("Exec failed: " + err.Error())
		os.Exit(1)
	}
}
