package create

import (
	"fmt"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createPodCmd = &cobra.Command{
	Use:   "pod <name> <image>",
	Short: "Create a new pod within a tenant",
	Long: `Creates a new pod. You need to specify the name and the image from which the pod should be created.
You can customize your deployment with the flags below or by using the interactive mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Create a pod...")
	},
}

func init() {
	createCmd.AddCommand(createPodCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	createCmd.Flags().BoolP("keep-alive", "", false, "Pod will be restarted upon termination.")
}
