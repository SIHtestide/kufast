package create

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"kufast/asyncOps"
	"os"
	"time"
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
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
		s.Prefix = "Creating Objects - Please wait!  "
		s.Start()

		res := asyncOps.NewPod(cmd, s, args)
		_ = <-res
		s.Stop()
		fmt.Println("Complete!")

	},
}

func init() {
	createCmd.AddCommand(createPodCmd)

	createPodCmd.Flags().BoolP("keep-alive", "", false, "Pod will be restarted upon termination.")
	createPodCmd.Flags().StringP("memory", "", "500Mi", "Limit the RAM usage for this namespace")
	createPodCmd.Flags().StringP("cpu", "", "500m", "Limit the CPU usage for this namespace")
	createPodCmd.Flags().StringP("target", "t", "", "The name of the node to deploy the pod")
	createPodCmd.Flags().StringP("deploy-secret", "d", "", "The name of the deployment secret to deploy this container.")
	createPodCmd.Flags().StringArrayP("secrets", "s", []string{}, "List of secret names to be introduced in the container as environment variables. The name equals the name of the secret")
}

func createPodInteractive(cmd *cobra.Command, args []string) []string {

	return args
}
