package delete

import (
	"fmt"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteNamespaceCmd = &cobra.Command{
	Use:   "delete namespace <namespace>",
	Short: "Delete the namespace including all users and pods in it.",
	Long: `Delete the namespace including all users and pods in it. This operation can only be executed by a cluster admin.
Please use with care! Deleted data cannot be restored.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("delete called")
	},
}

func init() {
	deleteCmd.AddCommand(deleteNamespaceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
