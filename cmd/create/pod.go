package create

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"kufast/clusterOperations"
	"kufast/tools"
)

// createCmd represents the create command
var createPodCmd = &cobra.Command{
	Use:   "pod <name> <image>",
	Short: "Create a new pod within a tenant-target",
	Long: `Creates a new pod. A pod is like a shell for a container in Kuebrnetes. 
You need to specify the name and the image from which the pod should be created.
You can customize your deployment with the flags below or by using the interactive mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		isInteractive, _ := cmd.Flags().GetBool("interactive")
		if isInteractive {
			args = createPodInteractive(cmd, args)
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

func init() {
	createCmd.AddCommand(createPodCmd)

	createPodCmd.Flags().BoolP("keep-alive", "", false, "Pod will be restarted upon termination.")
	createPodCmd.Flags().StringP("memory", "", "500Mi", "The amount of RAM the pod can use")
	createPodCmd.Flags().StringP("cpu", "", "500m", "The amount of CPU the pod can use")
	createPodCmd.Flags().StringP("storage", "", "1Gi", "The amount of storage the pod can use")
	createPodCmd.Flags().StringP("target", "", "", "The name of the node or group to deploy the pod to")
	createPodCmd.Flags().StringP("tenant", "", "", "The name of the tenant to be used to deploy the pod")
	createPodCmd.Flags().StringP("deploy-secret", "d", "", "The name of the deployment secret to deploy this container. This secret will be used to pull the image.")
	createPodCmd.Flags().StringArrayP("secrets", "s", []string{}, "List of secret names to be introduced in the container "+
		"as environment variables. The variable name will equal the name of the secret. Can be specified multiple times.")
	createPodCmd.Flags().Int32SliceP("port", "p", []int32{}, "A port the pod should expose. Can be specified multiple times.")
	createPodCmd.Flags().StringArrayP("cmd", "", []string{}, "An initial command to be issued at pod start. Required by a few containers.")
}

func createPodInteractive(cmd *cobra.Command, args []string) []string {
	fmt.Println("This routine will create a new pod for you. Please note that your credentials must be available" +
		"according to our Readme and you need to hand over secrets, ports and an init command through the command line arguments" +
		"if you intend to use them. More information is available under 'kufast get pod --help'")
	//No arguments provided
	if len(args) == 0 {
		args = append(args, tools.GetDialogAnswer("Please enter the name, the pod should have in Kubernetes."))
		args = append(args, tools.GetDialogAnswer("Please enter the dockerimage, the pod should have."))
	} else if len(args) != 2 {
		tools.HandleError(errors.New(tools.ERROR_WRONG_NUMBER_ARGUMENTS), cmd)
	}
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
