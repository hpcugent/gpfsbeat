package main

import (
	"os"

	"github.com/hpcugent/gpfsbeat/cmd"

	_ "github.com/hpcugent/gpfsbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
