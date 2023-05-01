package delete

import (
	"context"
	"errors"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kufast/tools"
	"os"
	"time"
)

// deleteCmd represents the delete command
var deleteNamespaceCmd = &cobra.Command{
	Use:   "namespace <namespace>",
	Short: "Delete the namespace including all users and pods in it.",
	Long: `Delete the namespace including all users and pods in it. This operation can only be executed by a cluster admin.
Please use with care! Deleted data cannot be restored.`,
	Run: func(cmd *cobra.Command, args []string) {

		//Check that exactly one arg has been provided (the namespace)
		if len(args) != 1 {
			tools.HandleError(errors.New("Too few or too many arguments provided."), cmd)
		}

		//Ensure user knows what he does
		answer := tools.GetDialogAnswer("Namespace " + args[0] + " will be deleted together with all users and pods! Continue? (No/yes)")
		if answer == "yes" {
			clientset, _, err := tools.GetUserClient(cmd)
			if err != nil {
				tools.HandleError(err, cmd)
			}

			err = clientset.CoreV1().Namespaces().Delete(context.TODO(), args[0], v1.DeleteOptions{})
			if err != nil {
				tools.HandleError(err, cmd)
			}

			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
			s.Prefix = "Deleting Objects - Please wait!  "
			s.Start()

			for true {
				_, err := clientset.CoreV1().Namespaces().Get(context.TODO(), args[0], v1.GetOptions{})
				if err != nil {
					break
				}
			}
			s.Stop()
			fmt.Println("Successfully deleted!")

		}
	},
}

func init() {
	deleteCmd.AddCommand(deleteNamespaceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
