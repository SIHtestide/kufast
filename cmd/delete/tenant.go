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

// deleteTenantCmd represents the delete tenant command
var deleteTenantCmd = &cobra.Command{
	Use:   "tenant <tenant>..",
	Short: "Delete tenants, their tenant-targets, pods, secrets and their credentials.",
	Long: `Delete tenants, their tenant-targets and their credentials. This operation can only be executed by a cluster admin.
Please use with care! Deleted data cannot be restored.`,
	Run: func(cmd *cobra.Command, args []string) {
		//Check that at least one tenant has been provided
		if len(args) < 1 {
			tools.HandleError(errors.New(tools.ERROR_WRONG_NUMBER_ARGUMENTS), cmd)
		}

		//Ensure user knows what he does
		answer := tools.GetDialogAnswer("Tenant will be deleted along with all deployment-targets, continue(yes/No)?")
		if answer == "yes" {

			s := tools.CreateStandardSpinner(tools.MESSAGE_DELETE_OBJECTS)

			for _, tenantName := range args {

				tenantTargets, err := clusterOperations.ListTargetsFromString(cmd, tenantName, false)
				if err != nil {
					tools.HandleError(err, cmd)
				}

				var deleteTargetOps []<-chan string
				var targetResults []string

				for _, tenantTarget := range tenantTargets {

					deleteTargetOps = append(deleteTargetOps, clusterOperations.DeleteTenantTarget(tenantTarget.Name, tenantName, cmd))
				}

				//Ensure all operations are done
				for _, op := range deleteTargetOps {
					targetResults = append(targetResults, <-op)
				}

				errorInDeletion := false
				for _, res := range targetResults {
					if res != "" {
						s.Stop()
						fmt.Println(res)
						s.Start()
						errorInDeletion = true
					}
				}
				if errorInDeletion {
					continue
				}

				err = clusterOperations.DeleteTenant(tenantName, cmd)
				if err != nil {
					tools.HandleError(err, cmd)
				}
			}

			s.Stop()
			fmt.Println(tools.MESSAGE_DONE)

		}
	},
}

// init is a helper function from cobra to initialize the command. It sets all flags, standard values and documentation for this command.
func init() {
	deleteCmd.AddCommand(deleteTenantCmd)

}
