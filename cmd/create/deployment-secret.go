package create

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"kufast/clusterOperations"
	"kufast/tools"
)

// createCmd represents the create command
var createDeploySecretCmd = &cobra.Command{
	Use:   "deploy-secret name",
	Short: "Create a new environment secret in this namespace",
	Long: `This command creates a new user and adds him to a namespace. You can select the namespace of the user.
Upon completion, the command yields the users credentials. This command will fail, if you do not have admin rights 
on the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			tools.HandleError(errors.New(tools.ERROR_WRONG_NUMBER_ARGUMENTS), cmd)
		}

		s := tools.CreateStandardSpinner(tools.MESSAGE_CREATE_OBJECTS)

		err := clusterOperations.CreateDeploymentSecret(args[0], cmd)
		if err != nil {
			s.Stop()
			tools.HandleError(err, cmd)
		}

		s.Stop()
		fmt.Println(tools.MESSAGE_DONE)

	},
}

func init() {
	createCmd.AddCommand(createDeploySecretCmd)

	createDeploySecretCmd.Flags().StringP("input", "", "", "Path to your .dockerfile to read credentials from.")
	createDeploySecretCmd.Flags().StringP("target", "", "", "The target for the secret (Needs to be the same as the pod using it.")
	createDeploySecretCmd.Flags().StringP("tenant", "", "", "The tenant owning the secret")
	createDeploySecretCmd.MarkFlagRequired("input")
	createDeploySecretCmd.MarkFlagFilename("input", "json")

}
