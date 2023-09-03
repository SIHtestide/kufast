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
package delete

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"kufast/clusterOperations"
	"kufast/tools"
)

// deletePodCmd represents the delete pod command
var deletePodCmd = &cobra.Command{
	Use:   "pod <pods>..",
	Short: "Delete the selected pod.",
	Long:  `Delete the selected pod including its storage. Please use with care! Deleted data cannot be restored.`,
	Run: func(cmd *cobra.Command, args []string) {

		//Check that exactly one arg has been provided (the namespace)
		if len(args) != 1 {
			tools.HandleError(errors.New(tools.ERROR_WRONG_NUMBER_ARGUMENTS), cmd)
		}

		//Ensure user knows what he does
		answer := tools.GetDialogAnswer("Pod " + args[0] + " will be deleted together with its storage and logs! Continue? (No/yes)")
		if answer == "yes" {

			s := tools.CreateStandardSpinner(tools.MESSAGE_DELETE_OBJECTS)

			var deleteTargetOps []<-chan string
			var targetResults []string

			for _, podName := range args {
				deleteTargetOps = append(deleteTargetOps, clusterOperations.DeletePod(cmd, podName))

			}

			//Ensure all operations are done
			for _, op := range deleteTargetOps {
				targetResults = append(targetResults, <-op)
			}

			for _, res := range targetResults {
				if res != "" {
					s.Stop()
					fmt.Println(res)
					s.Start()
				}
			}

			s.Stop()
			fmt.Println(tools.MESSAGE_DONE)

		}

	},
}

// init is a helper function from cobra to initialize the command. It sets all flags, standard values and documentation for this command.
func init() {
	deleteCmd.AddCommand(deletePodCmd)

	deletePodCmd.Flags().StringP("target", "", "", tools.DOCU_FLAG_TARGET)
	deletePodCmd.Flags().StringP("tenant", "", "", tools.DOCU_FLAG_TENANT)

}
