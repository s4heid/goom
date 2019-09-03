package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"net/url"
	"os"

	au "github.com/logrusorgru/aurora"
	"github.com/pkg/errors"
	. "github.com/s4heid/goom/config"
	"github.com/spf13/cobra"
)

func NewOpenCmd(configManager ConfigManager, browser Browser) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "open",
		Short:   "Open url in web browser",
		Long:    "Open url matching a configured alias in the default web browser",
		Aliases: []string{"o"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("invalid number of args specified")
			}

			config, err := configManager.ReadConfig()
			if err != nil {
				return errors.Wrap(err, "reading config")
			}

			if !containsAlias(config.Rooms, args[0]) {
				return fmt.Errorf("invalid alias %q", args[0])
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			config, err := configManager.ReadConfig()
			if err != nil {
				fmt.Println(au.Red(errors.Wrap(err, "reading config")))
				os.Exit(1)
			}

			url, err := createURL(config, args[0])
			if err != nil {
				fmt.Println(au.Red(errors.Wrap(err, "creating url")))
				os.Exit(1)
			}

			err = openURL(url, browser)
			if err != nil {
				fmt.Println(au.Red(fmt.Sprintf("cannot open url in browser: %v", err)))
				os.Exit(1)
			}

			fmt.Printf("Opening %q in the browser...", au.Green(url))
		},
	}
	return cmd
}

func createURL(config Config, alias string) (string, error) {
	roomIndex := len(config.Rooms)
	for r, room := range config.Rooms {
		if room.Alias == alias {
			roomIndex = r
			break
		}
	}
	if roomIndex == len(config.Rooms) {
		return "", errors.New("finding room in config")
	}

	tmpl, err := template.New("url").Parse(config.Url)
	if err != nil {
		return "", errors.Wrap(err, "creating template")
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, config.Rooms[roomIndex])
	if err != nil {
		return "", errors.Wrap(err, "executing template")
	}

	u, err := url.Parse(buf.String())
	if err != nil {
		return "", errors.Wrap(err, "parsing url")
	}

	return u.String(), nil
}

func containsAlias(rooms []Room, alias string) bool {
	for _, n := range rooms {
		if alias == n.Alias {
			return true
		}
	}
	return false
}

func init() {
	openCmd := NewOpenCmd(ConfigReader, DefaultBrowser{})
	rootCmd.AddCommand(openCmd)
}
