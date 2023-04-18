package create

import (
	"github.com/spf13/cobra"
	"kufast/trackerFactory"
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

		objectsCreated := 1
		pw := trackerFactory.CreateProgressWriter(objectsCreated)
		go trackerFactory.NewCreatePodTracker(cmd, pw, args)

		trackerFactory.HandleTracking(pw, objectsCreated)

	},
}

func init() {
	createCmd.AddCommand(createPodCmd)

	createPodCmd.Flags().BoolP("keep-alive", "", false, "Pod will be restarted upon termination.")
	createPodCmd.Flags().StringP("limit-memory", "m", "500Mi", "Limit the RAM usage for this namespace")
	createPodCmd.Flags().StringP("limit-cpu", "c", "500m", "Limit the CPU usage for this namespace")
	createPodCmd.Flags().StringP("target", "t", "", "The name of the node to deploy the pod")
	createPodCmd.Flags().StringArrayP("secrets", "s", []string{}, "List of secret names to be introduced in the container as environment variables. The name equals the name of the secret")
}

func createPodInteractive(cmd *cobra.Command, args []string) []string {

	return args
}
