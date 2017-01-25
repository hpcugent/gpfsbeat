package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/itkovian/gpfsbeat/beater"
)

func main() {
	err := beat.Run("gpfsbeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
