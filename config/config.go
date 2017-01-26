// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

// Config items
type Config struct {
	Period     time.Duration `config:"period"`
	Filesystem string        `config:"filesystem"`
}

// DefaultConfig should be overridden
var DefaultConfig = Config{
	Period:     1 * time.Second,
	Filesystem: "/",
}
