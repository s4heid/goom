package cmd

import (
	. "github.com/s4heid/goom/config"
)

//go:generate counterfeiter -o ./fakes/open.go . ConfigManager
type ConfigManager interface {
	ReadConfig() (Config, error)
}
