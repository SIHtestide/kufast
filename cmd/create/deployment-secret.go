/*
MIT License

Copyright (c) 2023 Stefan Pawlowski

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package create

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"kufast/clusterOperations"
	"kufast/tools"
)

// createDeploySecretCmd represents the create deploy-secret command
var createDeploySecretCmd = &cobra.Command{
	Use:   "deploy-secret name",
	Short: "Creates a deploy-secret in the specified target.",
	Long: `This command created a deploy-secret on the specified target. 
Deploy-secrets are required to create pods from private registries. This kind of secrets will be created
from dockerconfig files. A prerequisite to this command is, to store this file on your computer 
(default location is ~/.docker/config.json).
`,
	Run: func(cmd *cobra.Command, args []string) {
		isInteractive, _ := cmd.Flags().GetBool("interactive")
		if isInteractive {
			args = createDeploySecretInteractive(cmd)
		}

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

// createDeploySecretInteractive is a helper function to create a deploy-secret interactively
func createDeploySecretInteractive(cmd *cobra.Command) []string {
	fmt.Println(tools.MESSAGE_INTERACTIVE_IGNORE_INPUT)
	fmt.Println(`Please note, you need to specify the path to your dockercongigjson upfront.
	Consult the help page by calling this command with the --help flag for further information.`)
	var args []string
	args = append(args, tools.GetDialogAnswer("Please enter the name, the deploy secret should have in the system."))

	target := tools.GetDialogAnswer("Please select the tenant-target for your request. Leave empty for default. You can use 'kufast list tenant-targets' to get a list.")
	_ = cmd.Flags().Set("target", target)
	return args
}

// init is a helper function from cobra to initialize the command. It sets all flags, standard values and documentation for this command.
func init() {
	createCmd.AddCommand(createDeploySecretCmd)

	//target and tenant for the operation can be specified
	createDeploySecretCmd.Flags().StringP("target", "", "", tools.DOCU_FLAG_TARGET+" (Needs to be the same as the pod using it).")
	createDeploySecretCmd.Flags().StringP("tenant", "", "", tools.DOCU_FLAG_TENANT)

	//Input is required to get the secret
	createDeploySecretCmd.Flags().StringP("input", "", "", "Path to your .dockerconfigfile to read credentials from.")
	createDeploySecretCmd.MarkFlagRequired("input")
	createDeploySecretCmd.MarkFlagFilename("input", "json")

}
