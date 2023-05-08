package delete

import (
	"context"
	"errors"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kufast/clusterOperations"
	"kufast/tools"
	"os"
	"time"
)

// deleteTenantTargetCmd represents the delete tenant-target command
var deleteTenantTargetCmd = &cobra.Command{
	Use:   "tenant-target <tenant-target>",
	Short: "Delete the namespace including all users and pods in it.",
	Long: `Delete the namespace including all users and pods in it. This operation can only be executed by a cluster admin.
Please use with care! Deleted data cannot be restored.`,
	Run: func(cmd *cobra.Command, args []string) {

		//Ensure user knows what he does
		answer := tools.GetDialogAnswer("Namespaces will be deleted together with all users and pods! Continue? (No/yes)")
		if answer == "yes" {

			//Configblock
			clientset, _, err := tools.GetUserClient(cmd)
			if err != nil {
				fmt.Println(err.Error())
			}

			//Check that exactly one arg has been provided (the namespace)
			if len(args) != 1 {
				tools.HandleError(errors.New("Too few or too many arguments provided."), cmd)
			}

			tenantName, err := cmd.Flags().GetString("tenant")
			if err != nil {
				tools.HandleError(err, cmd)
			}

			//Activate spinner
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
			s.Prefix = "Deleting Objects - Please wait!  "
			s.Start()

			user, err := clientset.CoreV1().ServiceAccounts("default").Get(context.TODO(), tenantName+"-user", metav1.GetOptions{})
			if err != nil {
				s.Stop()
				tools.HandleError(err, cmd)
			}

			var results []int

			for _, tenantTargetName := range args {

				err = clientset.CoreV1().Namespaces().Delete(context.TODO(), tenantName+"-"+tenantTargetName, v1.DeleteOptions{})
				if err != nil {
					s.Stop()
					fmt.Println(err)
					s.Start()
					results = append(results, 1)
				} else {
					results = append(results, 0)
				}
			}

			//Remove capability from user
			for i, res := range results {
				if res == 0 {
					err := clusterOperations.DeleteTargetFromTenant(args[i], tenantName, cmd)
					if err != nil {
						s.Stop()
						tools.HandleError(err, cmd)
					}
				}
			}

			_, err = clientset.CoreV1().ServiceAccounts("default").Update(context.TODO(), user, metav1.UpdateOptions{})
			if err != nil {
				s.Stop()
				fmt.Println(err)
				s.Start()
			}
			s.Stop()
			fmt.Println(tools.MESSAGE_DONE)

		}

	},
}

// init is a helper function from cobra to initialize the command. It sets all flags, standard values and documentation for this command.
func init() {
	deleteCmd.AddCommand(deleteTenantTargetCmd)

	deleteTenantTargetCmd.Flags().StringP("tenant", "t", "", "The tenant owning this tenant-target.")
	_ = deleteTenantTargetCmd.MarkFlagRequired("tenant")

}
