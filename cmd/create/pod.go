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

// createPodCmd represents the create pod command
var createPodCmds = &cobra.Command{
	Use:   "pod <name> <image>",
	Short: "Create a new pod within a tenant-target",
	Long: `Creates a new pod within a tenant-target. A pod is like a shell for a container in Kuebrnetes. 
You need to specify the name and the image from which the pod should be created.
You can customize your deployment with the flags below or by using the interactive mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		isInteractive, _ := cmd.Flags().GetBool("interactive")
		if isInteractive {
			args = createPodInteractive(cmd)
		}
		if len(args) != 2 {
			tools.HandleError(errors.New(tools.ERROR_WRONG_NUMBER_ARGUMENTS), cmd)
		}
		s := tools.CreateStandardSpinner(tools.MESSAGE_CREATE_OBJECTS)

		res := clusterOperations.CreatePod(cmd, args)
		err := <-res
		s.Stop()
		if err != "" {
			tools.HandleError(errors.New(err), cmd)
		}
		fmt.Println(tools.MESSAGE_DONE)

	},
}

// createPodInteractive is a helper function to create a pod interactively
func createPodInteractive(cmd *cobra.Command) []string {
	fmt.Println(tools.MESSAGE_INTERACTIVE_IGNORE_INPUT)
	fmt.Println("This routine will create a new pod for you. Please note that your credentials must be available" +
		"according to our Readme and you need to hand over secrets, ports and an init command through the command line arguments" +
		"if you intend to use them. More information is available under 'kufast get pod --help'")
	//No arguments provided
	var args []string
	args = append(args, tools.GetDialogAnswer("Please enter the name, the pod should have in Kubernetes."))
	args = append(args, tools.GetDialogAnswer("Please enter the dockerimage, the pod should have."))

	limitCpu := tools.GetDialogAnswer("Please select the CPU limit for your pod, e.g. 250m (0.25 CPUs)")
	_ = cmd.Flags().Set("cpu", limitCpu)
	limitRam := tools.GetDialogAnswer("Please select the memory limit for your pod, e.g. 500Mi (500MB)")
	_ = cmd.Flags().Set("memory", limitRam)
	limitStorage := tools.GetDialogAnswer("Please select the storage limit for your pod, e.g. 2Gi (2GB)")
	_ = cmd.Flags().Set("storage", limitStorage)
	target := tools.GetDialogAnswer("Please select the tenant-target for your request. Leave empty for default. You can use 'kufast list tenant-targets' to get a list.")
	_ = cmd.Flags().Set("target", target)
	keepAlive := tools.GetDialogAnswer("Should we restart the container for you, if it completes? (y/N)")
	if keepAlive == "y" {
		_ = cmd.Flags().Set("keep-alive", "true")
	}
	deploySecret := tools.GetDialogAnswer("If your deployment needs a deploy secret, please enter it now. Leave empty, if no deploy-secret is required. Please note that the deploy-secret must be in the same tenant-target.")
	_ = cmd.Flags().Set("deploy-secret", deploySecret)

	return args
}

// init is a helper function from cobra to initialize the command. It sets all flags, standard values and documentation for this command.
func init() {
	createCmds.AddCommand(createPodCmds)

	//Settings for the pod
	createPodCmds.Flags().BoolP("keep-alive", "", false, "Pod will be restarted upon termination.")
	createPodCmds.Flags().StringP("memory", "", "500Mi", "The amount of RAM the pod can use")
	createPodCmds.Flags().StringP("cpu", "", "500m", "The amount of CPU the pod can use")
	createPodCmds.Flags().StringP("storage", "", "1Gi", "The amount of storage the pod can use")
	createPodCmds.Flags().StringP("deploy-secret", "d", "", "The name of the deployment secret to deploy this container. This secret will be used to pull the image.")
	createPodCmds.Flags().StringArrayP("secrets", "s", []string{}, "List of secret names to be introduced in the container "+
		"as environment variables. The variable name will equal the name of the secret. Can be specified multiple times.")
	createPodCmds.Flags().Int32SliceP("port", "p", []int32{}, "A port the pod should expose. Can be specified multiple times.")
	createPodCmds.Flags().StringArrayP("cmd", "", []string{}, "An initial command to be issued at pod start. Required by a few containers.")

	createPodCmds.Flags().StringP("target", "", "", tools.DOCU_FLAG_TARGET)
	createPodCmds.Flags().StringP("tenant", "", "", tools.DOCU_FLAG_TENANT)
}
