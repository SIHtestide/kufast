/*
MIT License

Copyright (c) 2023 Stefan Pawlowski

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
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
