package cmd

import (
	"github.com/s4heid/goom/config"
)

// ConfigManager initializes config if not set and reads it.
//go:generate counterfeiter -o ./fakes/open.go . ConfigManager
type ConfigManager interface {
	ReadConfig() (config.Config, error)
}
