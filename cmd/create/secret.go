package create

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"kufast/clusterOperations"
	"kufast/tools"
)

// createCmd represents the create command
var createSecretCmd = &cobra.Command{
	Use:   "secret name",
	Short: "Create a new environment secret in this namespace",
	Long: `This command creates a new user and adds him to a namespace. You can select the namespace of the user.
Upon completion, the command yields the users credentials. This command will fail, if you do not have admin rights 
on the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			tools.HandleError(errors.New(tools.ERROR_WRONG_NUMBER_ARGUMENTS), cmd)
		}

		//Get the secret
		secretData := tools.GetPasswordAnswer("Enter your secret here:")

		s := tools.CreateStandardSpinner(tools.MESSAGE_CREATE_OBJECTS)

		err := clusterOperations.CreateSecret(args[0], secretData, cmd)
		if err != nil {
			s.Stop()
			tools.HandleError(err, cmd)
		}

		s.Stop()
		fmt.Println(tools.MESSAGE_DONE)

	},
}

func init() {
	createCmd.AddCommand(createSecretCmd)

	createSecretCmd.Flags().StringP("target", "", "", "The target for the secret (Needs to be the same as the pod using it.")
	createSecretCmd.Flags().StringP("tenant", "", "", "The tenant owning the secret")

}
