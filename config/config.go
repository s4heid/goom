package config

import (
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Reader wraps an instance of viper
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

// ReadConfig reads a config file from viper and returns a Config object.
func (r *Reader) ReadConfig() (Config, error) {
	var c Config

	if err := r.ViperConfig.ReadInConfig(); err != nil {
		return c, errors.Wrap(err, "read config from viper")
	}

	err := r.ViperConfig.Unmarshal(&c)
	if err != nil {
		return c, errors.Wrap(err, "unmarshaling")
	}

	return c, nil
}

// SetConfig sets the config from the command line option. If the flag is
// not set, it defaults to filename .goom in the home directory
func (r *Reader) SetConfig(configPath string) error {
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

	return nil
}

// FindRoom returns the first room whose alias property matches a passed string.
func FindRoom(rooms []Room, alias string) (Room, error) {
	var r Room
	for _, r := range rooms {
		if alias == r.Alias {
			return r, nil
		}
	}
	return r, errors.New("finding room")
}
