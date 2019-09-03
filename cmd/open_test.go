package cmd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/s4heid/goom/cmd"
	fakes "github.com/s4heid/goom/cmd/fakes"
	. "github.com/s4heid/goom/config"
	"github.com/spf13/cobra"
)

var _ = Describe("open", func() {
	var (
		fm      *fakes.FakeConfigManager
		fb      *fakes.FakeBrowser
		execute = func(command *cobra.Command, args []string) error {
			command.SetArgs(args)
			command.SilenceUsage = true
			command.SilenceErrors = true
			return command.Execute()
		}
		testConfig Config = Config{
			Url: "https://potatoe/{{.Id}}",
			Rooms: []Room{
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

	It("succeeds when alias is in config", func() {
		err := execute(NewOpenCmd(fm, fb), []string{"jd"})
		Ω(err).ShouldNot(HaveOccurred())
		Ω(fb.OpenURLCallCount()).Should(Equal(1))
	})

	It("errors when alias is not in config", func() {
		err := execute(NewOpenCmd(fm, fb), []string{"not-jd"})
		Ω(err).Should(HaveOccurred())
		Ω(err).Should(MatchError("invalid alias \"not-jd\""))
		Ω(fb.OpenURLCallCount()).Should(Equal(0))
	})

	It("errors when no args are given", func() {
		err := execute(NewOpenCmd(fm, fb), []string{})
		Ω(err).Should(HaveOccurred())
		Ω(err).Should(MatchError("invalid number of args specified"))
		Ω(fb.OpenURLCallCount()).Should(Equal(0))
	})

	It("errors when wrong number of args are given", func() {
		err := execute(NewOpenCmd(fm, fb), []string{"jd", "another-jd"})
		Ω(err).Should(HaveOccurred())
		Ω(err).Should(MatchError("invalid number of args specified"))
		Ω(fb.OpenURLCallCount()).Should(Equal(0))
	})
})
