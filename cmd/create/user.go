package create

import (
	"fmt"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createUserCmd = &cobra.Command{
	Use:   "user <namespace>",
	Short: "Create a new user for a tenant",
	Long: `This command creates a new user and adds him to a namespace. You can select the namespace of the user.
Upon completion, the command yields the users credentials. This command will fail, if you do not have admin rights 
on the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Create user...")
	},
}

func init() {
	createCmd.AddCommand(createUserCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
