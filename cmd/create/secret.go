// Package create contains the cmd functions for the creation of objects in kufast
package create

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"kufast/clusterOperations"
	"kufast/tools"
)

// createSecretCmd represents the create secret command
var createSecretCmd = &cobra.Command{
	Use:   "secret name",
	Short: "Create a new environment secret in this namespace",
	Long: `This command creates a new user and adds him to a namespace. You can select the namespace of the user.
Upon completion, the command yields the users credentials. This command will fail, if you do not have admin rights 
on the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		isInteractive, _ := cmd.Flags().GetBool("interactive")
		if isInteractive {
			args = createSecretInteractive(cmd)
		}

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

// createSecretInteractive is a helper function to create a secret interactively
func createSecretInteractive(cmd *cobra.Command) []string {
	fmt.Println(tools.MESSAGE_INTERACTIVE_IGNORE_INPUT)
	var args []string
	args = append(args, tools.GetDialogAnswer("Please enter the name, the secret should have in the system."))

	target := tools.GetDialogAnswer("Please select the tenant-target for your request. Leave empty for default. You can use 'kufast list tenant-targets' to get a list.")
	_ = cmd.Flags().Set("target", target)
	return args
}

// init is a helper function from cobra to initialize the command. It sets all flags, standard values and documentation for this command.
func init() {
	createCmd.AddCommand(createSecretCmd)

	createSecretCmd.Flags().StringP("target", "", "", tools.DOCU_FLAG_TARGET+" (Needs to be the same as the pod using it).")
	createSecretCmd.Flags().StringP("tenant", "", "", tools.DOCU_FLAG_TENANT)

}
