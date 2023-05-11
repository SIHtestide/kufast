package list

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"kufast/clusterOperations"
	"kufast/tools"
	"os"
)

// listTargetsCmd represents the list targets command
var listTargetsCmd = &cobra.Command{
	Use:   "targets",
	Short: "List possible targets for the current tenant.",
	Long:  `List possible targets for the current tenant. A target is a node or a group of nodes, a tenant can deploy to.`,
	Run: func(cmd *cobra.Command, args []string) {

		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			tools.HandleError(err, cmd)
		}

		s := tools.CreateStandardSpinner(tools.MESSAGE_GET_OBJECTS)

		targets, err := clusterOperations.ListTargetsFromCmd(cmd, all)

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"NAME", "Type"})
		for _, target := range targets {
			t.AppendRow(table.Row{target.Name, target.AccessType})
		}

		s.Stop()
		t.AppendSeparator()
		t.Render()
	},
}

// init is a helper function from cobra to initialize the command. It sets all flags, standard values and documentation for this command.
func init() {
	listCmd.AddCommand(listTargetsCmd)
	listTargetsCmd.PersistentFlags().BoolP("all", "a", false, "List all tenant targets available on the instance. Admin use only!")
	listTargetsCmd.Flags().StringP("tenant", "", "", tools.DOCU_FLAG_TENANT)
	listTargetsCmd.MarkFlagsMutuallyExclusive("all", "tenant")
}
