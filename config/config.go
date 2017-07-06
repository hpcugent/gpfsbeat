// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

// Config items
type Config struct {
	Period             time.Duration `config:"period"`
	Devices            []string      `config:"devices"`
	MMRepQuotaCommand  string        `config:"mmrepquota"`
	MMLsFsCommand      string        `config:"mmlsfs"`
	MMDfCommand        string        `config:"mmdf"`
	MMLsFilesetCommand string        `config:"mmlsfileset"`
	UserIDField        string        `config:"__hpc_ugent_huppel_tierietoemba_user_id_field"`
}

// DefaultConfig should be overridden
var DefaultConfig = Config{
	Period:             1 * time.Second,
	Devices:            []string{"all"},
	MMRepQuotaCommand:  "mmrepquota",
	MMLsFsCommand:      "mmlsfs",
	MMDfCommand:        "mmdf",
	MMLsFilesetCommand: "mmlsfileset",
	UserIDField:        "none",
}
