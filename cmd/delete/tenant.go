package delete

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"kufast/clusterOperations"
	"kufast/tools"
)

// deleteCmd represents the delete command
var deleteTenantCmd = &cobra.Command{
	Use:   "tenant <tenant>..",
	Short: "Delete tenants, their namespaces and their credentials.",
	Long: `Delete a user and his credentials. This operation can only be executed by a cluster admin.
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
					deleteTargetOps = append(deleteTargetOps, clusterOperations.DeleteTenantTarget(tenantTarget.Name, cmd))
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
					}
				}
				if errorInDeletion {
					continue
				}

				clusterOperations.DeleteTe
			}

			for _, user := range args {
				deleteOps = append(deleteOps, clusterOperations.DeleteTenant(tenant, cmd))
			}

			for _, op := range deleteOps {
				results = append(results, <-op)
			}

			s.Stop()
			fmt.Println("Done!")

		}
	},
}

func init() {
	deleteCmd.AddCommand(deleteTenantCmd)

}
