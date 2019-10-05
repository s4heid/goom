package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"net/url"

	au "github.com/logrusorgru/aurora"
	"github.com/pkg/errors"
	"github.com/s4heid/goom/config"
	"github.com/spf13/cobra"
)

// NewOpenCmd returns a new cobra open command.
func NewOpenCmd(configManager ConfigManager, browser Browser, ioStreams IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "open",
		Short:   "Open url in web browser",
		Long:    "Open url matching a configured alias in the default web browser",
		Aliases: []string{"o"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("invalid number of args specified")
			}

			c, err := configManager.ReadConfig()
			if err != nil {
				return errors.Wrap(err, "reading config")
			}

			room, err := config.FindRoom(c.Rooms, args[0])
			if err != nil {
				return fmt.Errorf("alias %q does not exist", args[0])
			}

			url, err := createURL(c.Url, room)
			if err != nil {
				return errors.Wrap(err, "creating url")
			}

			err = openURL(url, browser)
			if err != nil {
				return errors.Wrap(err, "opening url in browser")
			}

			fmt.Fprint(ioStreams.Stdout, fmt.Sprintf("Opening %q in the browser...", au.Green(url)))
			return nil
		},
	}
	return cmd
}

func createURL(urlTemplate string, room config.Room) (string, error) {
	tmpl, err := template.New("url").Parse(urlTemplate)
	if err != nil {
		return "", errors.Wrap(err, "creating template")
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, room)
	if err != nil {
		return "", errors.Wrap(err, "executing template")
	}

	u, err := url.Parse(buf.String())
	if err != nil {
		return "", errors.Wrap(err, "parsing url")
	}

	return u.String(), nil
}
