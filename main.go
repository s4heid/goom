package main

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"text/template"

	"github.com/spf13/viper"
	"gopkg.in/urfave/cli.v1"
)

type Config struct {
	Url   string `yaml:"url"`
	Rooms []Room `yaml:"rooms"`
}

type Room struct {
	ID    string `yaml:"id"`
	Name  string `yaml:"name"`
	Alias string `yaml:"alias"`
}

func openBrowser(url string) {
	var err error

	log.Printf("Opening %q in the browser", url)
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		log.Fatal(err)
	}
}

func createURL(baseURL string, room Room) (string, error) {
	tmpl, err := template.New("url").Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("creating template: %e", err)
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, room)
	if err != nil {
		return "", fmt.Errorf("executing template: %e", err)
	}

	u, err := url.Parse(buf.String())
	if err != nil {
		return "", fmt.Errorf("malformed url %s", err)
	}

	return u.String(), nil
}

func main() {
	vp := viper.New()
	vp.SetConfigType("yaml")
	vp.SetConfigName("config")
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	configPath := filepath.Join(usr.HomeDir, ".gozoom")
	vp.AddConfigPath(configPath)
	vp.SetDefault("url", "https://zoom.us/j/{{.ID}}")
	vp.SetDefault("room", []Room{})
	if err := vp.ReadInConfig(); err != nil {
		// https://github.com/spf13/viper/issues/433
		err := vp.WriteConfigAs(filepath.Join(configPath, "config.yml"))
		if err != nil {
			log.Fatalf("cannot create config file at %q: %s", configPath, err)
		}
	}

	C := Config{}
	err = vp.Unmarshal(&C)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	app := cli.NewApp()
	app.Version = "0.0.1"
	app.Name = "gozoom"
	app.Usage = "A command line tool to jump into zoom meetings"
	app.UsageText = "gozoom [OPTIONS]"

	app.Action = func(c *cli.Context) error {
		if c.NArg() != 1 {
			log.Fatalf("invalid number of args %d", c.NArg())
		}

		for _, p := range C.Rooms {
			if p.Alias == c.Args().Get(0) {

				url, err := createURL(C.Url, p)
				if err != nil {
					log.Fatal(err)
				}

				log.Printf("Joining zoom meeting of %q\n", p.Name)
				openBrowser(url)

				break
			}
		}

		return nil
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
