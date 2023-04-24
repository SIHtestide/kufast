package delete

import (
	"fmt"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deletePodCmd = &cobra.Command{
	Use:   "delete pod <pod>",
	Short: "Delete the selected pod including its storage.",
	Long:  `Delete the selected pod including its storage. Please use with care! Deleted data cannot be restored.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("delete called")
	},
}

func init() {
	deleteCmd.AddCommand(deletePodCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
