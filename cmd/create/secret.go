package create

import (
	"context"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kufast/objectFactory"
	"kufast/tools"
	"os"
	"time"
)

// createCmd represents the create command
var createSecretCmd = &cobra.Command{
	Use:   "secret name",
	Short: "Create a new environment secret in this namespace",
	Long: `This command creates a new user and adds him to a namespace. You can select the namespace of the user.
Upon completion, the command yields the users credentials. This command will fail, if you do not have admin rights 
on the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {

		//Default config block
		clientset, _, err := tools.GetUserClient(cmd)
		if err != nil {
			tools.HandleError(err, cmd)
		}

		//Get the namespace
		namespaceName := tools.GetTenantTargetFromCmd(cmd)

		//Get the secret
		secretData := tools.GetPasswordAnswer("Enter your secret here:")

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
		s.Prefix = "Creating Objects - Please wait!  "
		s.Start()

		//create secret object
		secretObject := objectFactory.NewSecret(namespaceName, args[0], secretData)

		//Push secret
		clientset.CoreV1().Secrets(namespaceName).Create(context.TODO(), secretObject, metav1.CreateOptions{})

		s.Stop()
		fmt.Println("Done!")

	},
}

func init() {
	createCmd.AddCommand(createSecretCmd)

	createSecretCmd.Flags().StringP("target", "", "", "The target for the secret (Needs to be the same as the pod using it.")
	createSecretCmd.Flags().StringP("tenant", "", "", "The tenant owning the secret")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
