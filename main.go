package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/hpcugent/gpfsbeat/beater"
)

func main() {
	err := beat.Run("gpfsbeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
