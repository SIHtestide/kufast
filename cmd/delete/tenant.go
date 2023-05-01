package delete

import (
	"errors"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"kufast/clusterOperations"
	"kufast/tools"
	"os"
	"time"
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
			tools.HandleError(errors.New("provide at least one tenant"), cmd)
		}

		//Ensure user knows what he does
		answer := tools.GetDialogAnswer("Tenant will be deleted along with all deployment-targets, continue(yes/No)?")
		if answer == "yes" {

			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
			s.Prefix = "Deleting Objects - Please wait!  "
			s.Start()

			var deleteOps []<-chan int32
			var results []int32

			for _, user := range args {
				deleteOps = append(deleteOps, clusterOperations.DeleteUser(user, cmd, s))
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
