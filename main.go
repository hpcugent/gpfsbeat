package main

import (
	"os"

	"github.com/itkovian/gpfsbeat/cmd"

	_ "github.com/itkovian/gpfsbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
