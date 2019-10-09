package cmd

import (
	"encoding/json"
	"fmt"
	"text/tabwriter"

	"github.com/pkg/errors"
	"github.com/s4heid/goom/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// NewShowCmd returns a new cobra show command.
func NewShowCmd(configManager ConfigManager, ioStreams IOStreams) *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:     "show",
		Short:   "Show information of an url alias",
		Long:    "Show information of an url alias",
		Aliases: []string{"s"},
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

			var output []byte
			switch outputFormat {
			case "json":
				if output, err = json.MarshalIndent(room, "", "  "); err != nil {
					return err
				}
				fmt.Fprint(ioStreams.Stdout, string(output))
			case "yaml":
				if output, err = yaml.Marshal(room); err != nil {
					return err
				}
				fmt.Fprint(ioStreams.Stdout, "---\n"+string(output))
			default:
				if err = printTable(ioStreams, room); err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "output format in yaml or json")

	return cmd
}

func printTable(ioStreams IOStreams, room config.Room) error {
	w := tabwriter.NewWriter(ioStreams.Stdout, 0, 0, 3, ' ', tabwriter.TabIndent)
	fmt.Fprintf(w, "\nAlias\tName\tId\n%s\t%s\t%s\n", room.Alias, room.Name, room.Id)
	return w.Flush()
}
