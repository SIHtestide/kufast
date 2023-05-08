package delete

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"kufast/clusterOperations"
	"kufast/tools"
)

// deleteSecretCmd represents the delete secret command
var deleteSecretCmd = &cobra.Command{
	Use:   "secret <secret>..",
	Short: "Deletes a secret from a tenant-target.",
	Long: `Deletes a secret from a tenant-target.
Please use with care! Deleted data cannot be restored. Can be used for normal secrets and deploy-secrets.`,
	Run: func(cmd *cobra.Command, args []string) {
		//Check that exactly one arg has been provided (the namespace)
		if len(args) < 1 {
			tools.HandleError(errors.New(tools.ERROR_WRONG_NUMBER_ARGUMENTS), cmd)
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

// init is a helper function from cobra to initialize the command. It sets all flags, standard values and documentation for this command.
func init() {
	deleteCmd.AddCommand(deleteSecretCmd)

	deleteSecretCmd.Flags().StringP("target", "", "", tools.DOCU_FLAG_TARGET)
	deleteSecretCmd.Flags().StringP("tenant", "", "", tools.DOCU_FLAG_TENANT)

}
