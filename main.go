package main

import (
	"os"

	"github.com/hpcugent/gpfsbeat/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
