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
	"context"
	"errors"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kufast/clusterOperations"
	"kufast/tools"
	"os"
	"time"
)

// deleteTenantTargetCmd represents the delete tenant-target command
var deleteTenantTargetCmd = &cobra.Command{
	Use:   "tenant-target <tenant-target>",
	Short: "Delete tenant-targets of a tenant including all pods and secrets in it.",
	Long: `Delete tenant-targets of a tenant including all pods and secrets in it. This operation can only be executed by a cluster admin.
Please use with care! Deleted data cannot be restored.`,
	Run: func(cmd *cobra.Command, args []string) {

		//Ensure user knows what he does
		answer := tools.GetDialogAnswer("Namespaces will be deleted together with all users and pods! Continue? (No/yes)")
		if answer == "yes" {

			//Configblock
			clientset, _, err := tools.GetUserClient(cmd)
			if err != nil {
				fmt.Println(err.Error())
			}

			//Check that exactly one arg has been provided (the namespace)
			if len(args) != 1 {
				tools.HandleError(errors.New(tools.ERROR_WRONG_NUMBER_ARGUMENTS), cmd)
			}

			tenantName, err := cmd.Flags().GetString("tenant")
			if err != nil {
				tools.HandleError(err, cmd)
			}

			//Activate spinner
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
			s.Prefix = tools.MESSAGE_DELETE_OBJECTS
			s.Start()

			var results []int

			for _, tenantTargetName := range args {

				err = clientset.CoreV1().Namespaces().Delete(context.TODO(), tenantName+"-"+tenantTargetName, v1.DeleteOptions{})
				if err != nil {
					s.Stop()
					fmt.Println(err)
					s.Start()
					results = append(results, 1)
				} else {
					results = append(results, 0)
				}
			}

			//Remove capability from user
			for i, res := range results {
				if res == 0 {
					err := clusterOperations.DeleteTargetFromTenant(args[i], tenantName, cmd)
					if err != nil {
						s.Stop()
						tools.HandleError(err, cmd)
					}
				}
			}

			s.Stop()
			fmt.Println(tools.MESSAGE_DONE)

		}

	},
}

// init is a helper function from cobra to initialize the command. It sets all flags, standard values and documentation for this command.
func init() {
	deleteCmd.AddCommand(deleteTenantTargetCmd)

	deleteTenantTargetCmd.Flags().StringP("tenant", "t", "", "The tenant owning this tenant-target.")
	_ = deleteTenantTargetCmd.MarkFlagRequired("tenant")

}
