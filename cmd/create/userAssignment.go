package create

import (
	"fmt"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createUserAssignmentCmd = &cobra.Command{
	Use:   "userassign <namespace>",
	Short: "Create a new user assignment to a tenant",
	Long: `This command creates a new user assignment for a tenent. You must specify the namespace.
This command will fail, if you do not have admin rights 
on the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Create user assignment...")
	},
}

func init() {
	createCmd.AddCommand(createUserAssignmentCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
