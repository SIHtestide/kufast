package delete

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"kufast/clusterOperations"
	"kufast/tools"
)

// deleteCmd represents the delete command
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

func init() {
	deleteCmd.AddCommand(deletePodCmd)

	deletePodCmd.Flags().StringP("target", "", "", "The name of the node to deploy the pod")
	deletePodCmd.Flags().StringP("tenant", "", "", "The name of the tenant to deploy the pod to")

}
