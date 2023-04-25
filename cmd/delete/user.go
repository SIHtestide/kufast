package delete

import (
	"errors"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"kufast/asyncOps"
	"kufast/tools"
	"os"
	"time"
)

// deleteCmd represents the delete command
var deleteUserCmd = &cobra.Command{
	Use:   "user <user>..",
	Short: "Delete a user and his credentials.",
	Long: `Delete a user and his credentials. This operation can only be executed by a cluster admin.
Please use with care! Deleted data cannot be restored.`,
	Run: func(cmd *cobra.Command, args []string) {
		//Check that exactly one arg has been provided (the namespace)
		if len(args) < 1 {
			tools.HandleError(errors.New("Too few arguments provided."), cmd)
		}

		//Ensure user knows what he does
		answer := tools.GetDialogAnswer("User " + args[0] + " will be deleted! Continue? (No/yes)")
		if answer == "yes" {

			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
			s.Prefix = "Creating Objects - Please wait!  "
			s.Start()

			var deleteOps []<-chan int32
			var results []int32

			for _, user := range args {
				deleteOps = append(deleteOps, asyncOps.DeleteUser(user, cmd, s))
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
	deleteCmd.AddCommand(deleteUserCmd)

}
