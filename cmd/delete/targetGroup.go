package delete

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"kufast/clusterOperations"
	"kufast/tools"
)

// deleteCmd represents the delete command
var deleteTargetGroupCmd = &cobra.Command{
	Use:   "targetGroup <targetgroup>..",
	Short: "Delete a user and his credentials.",
	Long: `Delete a user and his credentials. This operation can only be executed by a cluster admin.
Please use with care! Deleted data cannot be restored.`,
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
			fmt.Println("Done!")

		}
	},
}

func init() {
	deleteCmd.AddCommand(deleteTargetGroupCmd)

	deleteTargetGroupCmd.Flags().StringP("target", "", "", "The name of the node to deploy the pod")
	deleteTargetGroupCmd.Flags().StringP("tenant", "", "", "The name of the tenant to deploy the pod to")

}
