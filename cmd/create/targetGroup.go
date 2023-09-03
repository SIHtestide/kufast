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
package create

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"kufast/clusterOperations"
	"kufast/tools"
)

// createTargetGroupCmd represents the create target-group command
var createTargetGroupCmd = &cobra.Command{
	Use:   "target-group <name> <nodes>",
	Short: "Create a target-group within the cluster",
	Long: `This command creates a new target-group and assigns it to the specified nodes.
Target-groups can be used to define a tenant-target that can deploy to a group of nodes,
instead of a single node.`,
	Run: func(cmd *cobra.Command, args []string) {
		isInteractive, _ := cmd.Flags().GetBool("interactive")
		if isInteractive {
			args = createTargetGroupInteractive()
		}

		if len(args) < 2 {
			tools.HandleError(errors.New(tools.ERROR_WRONG_NUMBER_ARGUMENTS), cmd)
		}

		s := tools.CreateStandardSpinner(tools.MESSAGE_CREATE_OBJECTS)

		err := clusterOperations.SetTargetGroupToNodes(args[0], args[1:], cmd)
		if err != nil {
			tools.HandleError(err, cmd)
		}

		s.Stop()
		fmt.Println(tools.MESSAGE_DONE)

	},
}

// createTargetGroupInteractive is a helper function to create a target-group interactively
func createTargetGroupInteractive() []string {
	fmt.Println(tools.MESSAGE_INTERACTIVE_IGNORE_INPUT)
	var args []string
	args = append(args, tools.GetDialogAnswer("Please specify the name of the target-group. It has to be alphanumeric"))
	for true {
		args = append(args, tools.GetDialogAnswer("Please enter the name, of a node, that should be part of the target-group."))
		next := tools.GetDialogAnswer("Do you want to add another node? (y/N)")
		if next != "y" {
			break
		}
	}
	return args
}

// init is a helper function from cobra to initialize the command. It sets all flags, standard values and documentation for this command.
func init() {
	createCmd.AddCommand(createTargetGroupCmd)

}
