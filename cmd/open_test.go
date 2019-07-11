package cmd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/s4heid/goom/cmd"
	configfakes "github.com/s4heid/goom/cmd/fakes"
	. "github.com/s4heid/goom/config"
	"github.com/spf13/cobra"
)

var _ = Describe("open", func() {
	var (
		fm      *configfakes.FakeManager
		execute = func(command *cobra.Command, args []string) error {
			command.SetArgs(args)
			return command.Execute()
		}
		testConfig Config = Config{
			Url: "https://potatoe/{{.Id}}",
			Rooms: []Room{
				Room{
					Alias: "jd",
					Id:    "123",
					Name:  "John Doe",
				},
			},
		}
	)

	BeforeEach(func() {
		fm = &configfakes.FakeManager{}
		fm.ReadConfigReturns(testConfig, nil)
	})

	It("succeeds when alias is in config", func() {
		err := execute(NewOpenCmd(fm), []string{"jd"})
		Ω(err).ShouldNot(HaveOccurred())
	})

	It("errors when alias is not in config", func() {
		err := execute(NewOpenCmd(fm), []string{"not-jd"})
		Ω(err).Should(HaveOccurred())
		Ω(err).Should(MatchError("config does not contain alias \"not-jd\""))
	})

	It("errors when no args are given", func() {
		err := execute(NewOpenCmd(fm), []string{})
		Ω(err).Should(HaveOccurred())
		Ω(err).Should(MatchError("invalid number of args specified"))
	})

	It("errors when wrong number of args are given", func() {
		err := execute(NewOpenCmd(fm), []string{"jd", "another-jd"})
		Ω(err).Should(HaveOccurred())
		Ω(err).Should(MatchError("invalid number of args specified"))
	})
})
