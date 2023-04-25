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
var deletePodCmd = &cobra.Command{
	Use:   "pod <pod>",
	Short: "Delete the selected pod including its storage.",
	Long:  `Delete the selected pod including its storage. Please use with care! Deleted data cannot be restored.`,
	Run: func(cmd *cobra.Command, args []string) {

		//Check that exactly one arg has been provided (the namespace)
		if len(args) != 1 {
			tools.HandleError(errors.New("Too few or too many arguments provided."), cmd)
		}

		//Ensure user knows what he does
		answer := tools.GetDialogAnswer("Pod " + args[0] + " will be deleted together with its storage and logs! Continue? (No/yes)")
		if answer == "yes" {

			clientset, _, err := tools.GetUserClient(cmd)
			if err != nil {
				tools.HandleError(err, cmd)
			}

			namespaceName, _ := tools.GetNamespaceFromUserConfig(cmd)

			err = clientset.CoreV1().Pods(namespaceName).Delete(context.TODO(), args[0], v1.DeleteOptions{})
			if err != nil {
				tools.HandleError(err, cmd)
			}

			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
			s.Prefix = "Deleting Objects - Please wait!  "
			s.Start()

			//Check for the pod been deleted from the system
			for true {
				_, err := clientset.CoreV1().Pods(namespaceName).Get(context.TODO(), args[0], v1.GetOptions{})
				if err != nil {
					break
				}
			}
			s.Stop()
			fmt.Println("Done!")

		}

	},
}

func init() {
	deleteCmd.AddCommand(deletePodCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
