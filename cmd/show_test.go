package cmd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/s4heid/goom/cmd"
	fakes "github.com/s4heid/goom/cmd/fakes"
	goomconfig "github.com/s4heid/goom/config"
	"github.com/spf13/cobra"
)

var _ = Describe("NewShowCmd", func() {
	var (
		fm      *fakes.FakeConfigManager
		execute = func(command *cobra.Command, args []string) error {
			command.SetArgs(args)
			command.SilenceUsage = true
			command.SilenceErrors = true
			return command.Execute()
		}
		testConfig = goomconfig.Config{
			Url: "https://potatoe/{{.Id}}",
			Rooms: []goomconfig.Room{
				{
					Alias: "po",
					Id:    "123",
					Name:  "Potatoe",
				},
				{
					Alias: "be",
					Id:    "321",
					Name:  "Beer",
				},
			},
		}
	)

	BeforeEach(func() {
		fm = &fakes.FakeConfigManager{}
		fm.ReadConfigReturns(testConfig, nil)
	})

	It("Prints detail information in table format by default", func() {
		ioStreams, _, out, _ := cmd.NewTestIOStreams()
		err := execute(cmd.NewShowCmd(fm, ioStreams), []string{"po"})

		Ω(err).ShouldNot(HaveOccurred())
		Ω(out.String()).Should(Equal("\nAlias   Name      Id\npo      Potatoe   123\n"))
	})

	It("Prints output in table format when --output flag receives invalid input", func() {
		ioStreams, _, out, _ := cmd.NewTestIOStreams()
		err := execute(cmd.NewShowCmd(fm, ioStreams), []string{"--output=invalid", "po"})

		Ω(err).ShouldNot(HaveOccurred())
		Ω(out.String()).Should(Equal("\nAlias   Name      Id\npo      Potatoe   123\n"))
	})

	It("Prints output in json format when --json flag is specified", func() {
		ioStreams, _, out, _ := cmd.NewTestIOStreams()
		err := execute(cmd.NewShowCmd(fm, ioStreams), []string{"--output=json", "po"})

		Ω(err).ShouldNot(HaveOccurred())
		Ω(out.String()).Should(Equal("{\n  \"id\": \"123\",\n  \"name\": \"Potatoe\",\n  \"alias\": \"po\"\n}"))
	})

	It("Prints output in yaml format when --yaml flag is specified", func() {
		ioStreams, _, out, _ := cmd.NewTestIOStreams()
		err := execute(cmd.NewShowCmd(fm, ioStreams), []string{"--output=yaml", "be"})

		Ω(err).ShouldNot(HaveOccurred())
		Ω(out.String()).Should(Equal("---\nid: \"321\"\nname: Beer\nalias: be\n"))
	})
})
