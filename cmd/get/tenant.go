package get

import (
	"errors"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"kufast/clusterOperations"
	"kufast/tools"
	"os"
)

// getCmd represents the get command
var getTenantCmd = &cobra.Command{
	Use:   "secret <secret>",
	Short: "Gain information about a deployed pod.",
	Long:  `Gain information about a deployed pod.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			tools.HandleError(errors.New(tools.ERROR_WRONG_NUMBER_ARGUMENTS), cmd)
		}

		s := tools.CreateStandardSpinner(tools.MESSAGE_GET_OBJECTS)

		tenant, err := clusterOperations.GetTenantFromString(cmd, args[0])
		if err != nil {
			tools.HandleError(err, cmd)
		}

		targets, err := clusterOperations.ListTargetsFromString(cmd, args[0], false)
		if err != nil {
			tools.HandleError(err, cmd)
		}
		var groupTargets []string
		var nodeTargets []string

		for _, target := range targets {
			if target.AccessType == "node" {
				groupTargets = append(groupTargets, target.Name)
			} else {
				nodeTargets = append(nodeTargets, target.Name)
			}
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"ATTRIBUTE", "VALUE"})
		t.AppendRow(table.Row{"Name", tenant.Name})
		t.AppendRow(table.Row{"Node Access", nodeTargets})
		t.AppendRow(table.Row{"Group Access", groupTargets})

		s.Stop()

		t.AppendSeparator()
		t.Render()

	},
}

func init() {
	getCmd.AddCommand(getTenantCmd)

}
