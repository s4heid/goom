package config_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	goomconfig "github.com/s4heid/goom/config"

	"github.com/spf13/viper"
)

var subject *goomconfig.Reader
var configContent, configPath, tmpdir string

var _ = BeforeEach(func() {
	configContent = string(`
{
	"url": "https://my-fake-room/{{.Id}}",
	"rooms": [{
		"id": "1234",
		"name": "John Doe",
		"alias": "jd"
	}]
}`)
	tmpdir, _ = ioutil.TempDir("", "tests")
	configPath = filepath.Join(tmpdir, ".goom.json")
	_ = ioutil.WriteFile(configPath, []byte(configContent), 0666)
})

var _ = Describe("ReadConfig", func() {
	It("reads the config in json format", func() {
		subject = &goomconfig.Reader{
			ViperConfig: viper.New(),
		}
		subject.ViperConfig.SetConfigFile(configPath)

		config, err := subject.ReadConfig()

		Ω(err).ShouldNot(HaveOccurred())
		Ω(config.Url).Should(Equal("https://my-fake-room/{{.Id}}"))
		Ω(config.Rooms).Should(Equal(
			[]goomconfig.Room{
				{
					Alias: "jd",
					Id:    "1234",
					Name:  "John Doe",
				},
			},
		))
	})

	Context("when config is in yaml format", func() {
		BeforeEach(func() {
			configContent = string(`
url: https://my-fake-room/{{.Id}}
rooms:
  id: 4321
  name: Yoshi Doe
  alias: yd
`)
			configPath = filepath.Join(tmpdir, ".goom.yaml")
			_ = ioutil.WriteFile(configPath, []byte(configContent), 0666)
		})

		It("reads the config", func() {
			subject = &goomconfig.Reader{
				ViperConfig: viper.New(),
			}
			subject.ViperConfig.SetConfigFile(configPath)

			config, err := subject.ReadConfig()

			Ω(err).ShouldNot(HaveOccurred())
			Ω(subject.ViperConfig.ConfigFileUsed()).Should(Equal(configPath))
			Ω(config.Url).Should(Equal("https://my-fake-room/{{.Id}}"))
			Ω(config.Rooms).Should(Equal(
				[]goomconfig.Room{
					{
						Alias: "yd",
						Id:    "4321",
						Name:  "Yoshi Doe",
					},
				},
			))
		})
	})

	Context("When config filepath is empty string", func() {
		var originalHome string

		BeforeEach(func() {
			originalHome = os.Getenv("HOME")
			os.Setenv("HOME", tmpdir)
		})

		AfterEach(func() {
			os.RemoveAll(tmpdir)
			os.Setenv("HOME", originalHome)
		})

		It("reads the config from home directory if configpath is empty string", func() {
			subject = &goomconfig.Reader{
				ViperConfig: viper.New(),
			}
			subject.ViperConfig.SetConfigFile(configPath)

			_, err := subject.ReadConfig()

			Ω(subject.ViperConfig.ConfigFileUsed()).Should(Equal(configPath))
			Ω(err).Should(BeNil())
		})
	})
})

var _ = Describe("SetConfig", func() {
	It("sets config file name .goom when filepath is specified", func() {
		subject = &goomconfig.Reader{
			ViperConfig: viper.New(),
		}

		err := subject.SetConfig(configPath)

		Ω(err).ShouldNot(HaveOccurred())
		Ω(subject.ViperConfig.ConfigFileUsed()).Should(Equal(configPath))
	})
})
