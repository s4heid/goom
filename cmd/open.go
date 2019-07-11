package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	. "github.com/logrusorgru/aurora"
	"github.com/pkg/errors"
	. "github.com/s4heid/goom/config"
	"github.com/spf13/cobra"
)

//go:generate counterfeiter -o ./fakes/open.go . Manager
type Manager interface {
	ReadConfig() (Config, error)
}

func NewOpenCmd(reader Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open",
		Short: "Open url in web browser",
		Long:  "Open url matching a configured alias in the default web browser",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("invalid number of args specified")
			}

			config, err := reader.ReadConfig()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if !contains(config.Rooms, args[0]) {
				return fmt.Errorf("config does not contain alias %q", args[0])
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			config, err := reader.ReadConfig()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = Open(config, args...)
			if err != nil {
				fmt.Println(Red(fmt.Sprintf("cannot open room: %v", err)))
				os.Exit(1)
			}
		},
	}
	return cmd
}

func Open(config Config, args ...string) error {
	url, err := createURL(config, args[0])
	if err != nil {
		return errors.Wrap(err, "creating url")
	}

	err = openBrowser(url)
	if err != nil {
		return errors.Wrap(err, "opening url in browser")
	}

	return nil
}

func contains(rooms []Room, alias string) bool {
	for _, n := range rooms {
		if alias == n.Alias {
			return true
		}
	}
	return false
}

func openBrowser(url string) error {
	var err error

	fmt.Printf("Opening %q in the browser...", Green(url))
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command(
			filepath.Join(
				os.Getenv("SYSTEMROOT"),
				"System32", "rundll32.exe"),
			"url.dll,FileProtocolHandler",
			url,
		).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = errors.New("unsupported platform")
	}

	return err
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

func init() {
	openCmd := NewOpenCmd(ConfigReader)
	rootCmd.AddCommand(openCmd)
}
