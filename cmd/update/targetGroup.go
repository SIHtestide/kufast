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
package update

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"kufast/clusterOperations"
	"kufast/tools"
)

// updateTargetGroupCmd represents the update target-group command
var updateTargetGroupCmd = &cobra.Command{
	Use:   "target-group <name> <nodes>..",
	Short: "Update the nodes on an existing target group.",
	Long: `Update the nodes on an existing target group. Specify all nodes that should be in the group after the reassignment. 
 Already existing pods on nodes will not be affected of this change.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 2 {
			tools.HandleError(errors.New(tools.ERROR_WRONG_NUMBER_ARGUMENTS), cmd)
		}

		s := tools.CreateStandardSpinner(tools.MESSAGE_UPDATE_OBJECTS)

		if clusterOperations.IsValidTarget(cmd, args[0], true) {
			err := clusterOperations.SetTargetGroupToNodes(args[0], args[:1], cmd)
			if err != nil {
				tools.HandleError(err, cmd)
			}
		}

		s.Stop()
		fmt.Println(tools.MESSAGE_DONE)

	},
}

// init is a helper function from cobra to initialize the command. It sets all flags, standard values and documentation for this command.
func init() {
	updateCmd.AddCommand(updateTargetGroupCmd)

}
