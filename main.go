package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
	"gopkg.in/urfave/cli.v1"
)

type Config struct {
	Room   string   `yaml:"room"`
	People []People `yaml:"people"`
}

type People struct {
	ID    string `yaml:"id"`
	Name  string `yaml:"name"`
	Alias string `yaml:"alias"`
}

func openBrowser(url string) {
	var err error

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

func createURL(baseURL, id string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		fmt.Println("malformed url: ", err.Error())
		return "", err
	}
	u.Path = path.Join(u.Path, id)

	return u.String(), nil
}

func main() {
	vp := viper.New()
	vp.SetConfigType("yaml")
	vp.SetConfigName("config")
	// vp.AddConfigPath("$HOME/.gozoom/")
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	vp.AddConfigPath(filepath.Join(usr.HomeDir, ".gozoom"))
	if err := vp.ReadInConfig(); err != nil {
		panic(fmt.Errorf("failed to create %s", vp.ConfigFileUsed()))
	}

	C := Config{}
	err = vp.Unmarshal(&C)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	app := cli.NewApp()
	app.Action = func(c *cli.Context) error {
		if c.NArg() != 1 {
			panic(fmt.Errorf("invalid number of args %d", c.NArg()))
		}

		for _, p := range C.People {
			if p.Alias == c.Args().Get(0) {

				url, err := createURL(C.Room, p.ID)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf("Joining zoom meeting of %q", p.Name)
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
