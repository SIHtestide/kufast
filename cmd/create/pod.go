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
	Short: "Create a new pod within a tenant",
	Long: `Creates a new pod. You need to specify the name and the image from which the pod should be created.
You can customize your deployment with the flags below or by using the interactive mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		isInteractive, _ := cmd.Flags().GetBool("interactive")
		if isInteractive {
			args = createPodInteractive(cmd, args)
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
	createPodCmd.Flags().StringP("memory", "", "500Mi", "Limit the RAM usage for this namespace")
	createPodCmd.Flags().StringP("cpu", "", "500m", "Limit the CPU usage for this namespace")
	createPodCmd.Flags().StringP("storage", "", "1Gi", "The amount of storage the pod can use")
	createPodCmd.Flags().StringP("target", "", "", "The name of the node to deploy the pod")
	createPodCmd.Flags().StringP("tenant", "", "", "The name of the tenant to deploy the pod to")
	createPodCmd.Flags().StringP("deploy-secret", "d", "", "The name of the deployment secret to deploy this container.")
	createPodCmd.Flags().StringArrayP("secrets", "s", []string{}, "List of secret names to be introduced in the container as environment variables. The name equals the name of the secret")
}

func createPodInteractive(cmd *cobra.Command, args []string) []string {

	return args
}
