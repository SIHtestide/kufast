package delete

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"kufast/clusterOperations"
	"kufast/tools"
)

// deleteCmd represents the delete command
var deleteSecretCmd = &cobra.Command{
	Use:   "secret <secret>..",
	Short: "Delete a user and his credentials.",
	Long: `Delete a user and his credentials. This operation can only be executed by a cluster admin.
Please use with care! Deleted data cannot be restored.`,
	Run: func(cmd *cobra.Command, args []string) {
		//Check that exactly one arg has been provided (the namespace)
		if len(args) < 1 {
			tools.HandleError(errors.New("Too few arguments provided."), cmd)
		}

		//Ensure user knows what he does
		answer := tools.GetDialogAnswer("Secrets will be deleted! Continue? (No/yes)")
		if answer == "yes" {

			s := tools.CreateStandardSpinner(tools.MESSAGE_DELETE_OBJECTS)

			var deleteOps []<-chan string
			var results []string

			for _, secret := range args {
				deleteOps = append(deleteOps, clusterOperations.DeleteSecret(secret, cmd))
			}

			for _, op := range deleteOps {
				results = append(results, <-op)
			}

			for _, res := range results {
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

func init() {
	deleteCmd.AddCommand(deleteSecretCmd)

	deleteSecretCmd.Flags().StringP("target", "", "", "The name of the node to deploy the pod")
	deleteSecretCmd.Flags().StringP("tenant", "", "", "The name of the tenant to deploy the pod to")

}
