package config

import (
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

////go:generate counterfeiter -o ./fakes/config.go . Manager
//type Manager interface {
//	ReadConfig() (Config, error)
//}

type Reader struct {
	ViperConfig *viper.Viper
}

type Config struct {
	Url   string `yaml:"url" json:"url"`
	Rooms []Room `yaml:"rooms" json:"rooms"`
}

type Room struct {
	Id    string `yaml:"id" json:"id"`
	Name  string `yaml:"name" json:"name"`
	Alias string `yaml:"alias" json:"alias"`
}

func (r Reader) ReadConfig() (Config, error) {
	var c Config

	if err := r.ViperConfig.ReadInConfig(); err != nil {
		return c, errors.Wrap(err, "reading config")
	}

	err := r.ViperConfig.Unmarshal(&c)
	if err != nil {
		return c, errors.Wrap(err, "unmarshaling")
	}

	return c, nil
}

func (r Reader) InitConfig(configPath string) error {
	if configPath != "" {
		r.ViperConfig.SetConfigFile(configPath)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			return errors.Wrap(err, "locating homedir")
		}

		r.ViperConfig.AddConfigPath(home)
		r.ViperConfig.SetConfigName(".goom")
	}

	r.ViperConfig.AutomaticEnv()

	if err := r.ViperConfig.ReadInConfig(); err != nil {
		return errors.Wrap(err, "reading config")
	}

	return nil
}

// func Run(manager Manager) (Config, error) {
// 	config, err := manager.ReadConfig()
// 	return config, err
// }
