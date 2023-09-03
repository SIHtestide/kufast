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
package get

import (
	"errors"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"kufast/clusterOperations"
	"kufast/tools"
	"os"
)

// getTenantCmd represents the get tenant command
var getTenantCmd = &cobra.Command{
	Use:   "tenant <tenant name>",
	Short: "Gain information about a deployed tenant.",
	Long:  `Gain information about a deployed tenant. Output includes name, node access and group access`,
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
			if target.AccessType == "group" {
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

// init is a helper function from cobra to initialize the command. It sets all flags, standard values and documentation for this command.
func init() {
	getCmd.AddCommand(getTenantCmd)

}
