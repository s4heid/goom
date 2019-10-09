package cmd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/s4heid/goom/cmd"
	fakes "github.com/s4heid/goom/cmd/fakes"
	goomconfig "github.com/s4heid/goom/config"
	"github.com/spf13/cobra"
)

var _ = Describe("cmd.NewOpenCmd", func() {
	var (
		fm      *fakes.FakeConfigManager
		fb      *fakes.FakeBrowser
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
					Alias: "jd",
					Id:    "123",
					Name:  "John Doe",
				},
			},
		}
	)

	BeforeEach(func() {
		fm = &fakes.FakeConfigManager{}
		fb = &fakes.FakeBrowser{}
		fm.ReadConfigReturns(testConfig, nil)
		fb.OpenURLReturns(nil)
	})

	It("opens the associated url of an alias", func() {
		ioStreams, _, out, _ := cmd.NewTestIOStreams()
		err := execute(cmd.NewOpenCmd(fm, fb, ioStreams), []string{"jd"})
		Ω(err).ShouldNot(HaveOccurred())
		Ω(out.String()).Should(Equal("Opening \x1b[32m\"https://potatoe/123\"\x1b[0m (John Doe) in the browser..."))
		Ω(fb.OpenURLCallCount()).Should(Equal(1))
	})

	It("errors when alias is not in config", func() {
		ioStreams, _, _, _ := cmd.NewTestIOStreams()
		err := execute(cmd.NewOpenCmd(fm, fb, ioStreams), []string{"not-jd"})
		Ω(err).Should(MatchError("alias \"not-jd\" does not exist"))
		Ω(fb.OpenURLCallCount()).Should(Equal(0))
	})

	It("errors when wrong number of args are given", func() {
		ioStreams, _, _, _ := cmd.NewTestIOStreams()
		err := execute(cmd.NewOpenCmd(fm, fb, ioStreams), []string{"jd", "another-jd"})
		Ω(err).Should(MatchError("invalid number of args specified"))
		Ω(fb.OpenURLCallCount()).Should(Equal(0))
	})
})
