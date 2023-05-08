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
	Short: "Creates a deploy-secret in the specified target.",
	Long: `This command created a deploy-secret on the specified target. Deploy-secrets are required to create
`,
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
	createDeploySecretCmd.Flags().StringP("tenant", "", "", "The tenant for this operation")
	createDeploySecretCmd.MarkFlagRequired("input")
	createDeploySecretCmd.MarkFlagFilename("input", "json")

}
