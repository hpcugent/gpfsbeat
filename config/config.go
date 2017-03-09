// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

// Config items
type Config struct {
	Period            time.Duration `config:"period"`
	Device            []string      `config:"device"`
	MMRepQuotaCommand string        `config:"mmrepquota"`
	MMLsFsCommand     string        `config:"mmlsfs"`
	MMDfCommand       string        `config:"mmsf"`
}

// DefaultConfig should be overridden
var DefaultConfig = Config{
	Period:            1 * time.Second,
	Device:            []string{"all"},
	MMRepQuotaCommand: "mmrepquota",
	MMLsFsCommand:     "mmlsfs",
	MMDfCommand:       "mmdf",
}
