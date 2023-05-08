package delete

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"kufast/clusterOperations"
	"kufast/tools"
)

// deleteTargetGroupCmd represents the delete target-group command
var deleteTargetGroupCmd = &cobra.Command{
	Use:   "target-group <target-group>..",
	Short: "Deletes a target-group from the cluster.",
	Long: `Deletes a target-group from the cluster. This operation can only be executed by a cluster admin.
Please use with care! Tenant-targets pointing to these target-groups remain intact, but cannot deploy new pods.`,
	Run: func(cmd *cobra.Command, args []string) {
		//Check that exactly one arg has been provided (the namespace)
		if len(args) < 1 {
			tools.HandleError(errors.New("Too few arguments provided."), cmd)
		}

		//Ensure user knows what he does
		answer := tools.GetDialogAnswer("Targetgroup will be deleted! Spaces with that target group remain intact but are unable to deploy! Continue (yes/No)")
		if answer == "yes" {

			s := tools.CreateStandardSpinner(tools.MESSAGE_DELETE_OBJECTS)

			for _, group := range args {
				err := clusterOperations.DeleteTargetGroupFromNodes(group, cmd)
				if err != nil {
					s.Stop()
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
	deleteCmd.AddCommand(deleteTargetGroupCmd)

	deleteTargetGroupCmd.Flags().StringP("target", "", "", tools.DOCU_FLAG_TARGET)
	deleteTargetGroupCmd.Flags().StringP("tenant", "", "", tools.DOCU_FLAG_TENANT)

}
