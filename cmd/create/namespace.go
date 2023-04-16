package create

import (
	"fmt"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createNamespaceCmd = &cobra.Command{
	Use:   "namespace <name>",
	Short: "Create a new namespace for a tenant",
	Long: `This command creates a new namespace for a tenant. You can select the name and set limits to the namespace.
This command will fail, if you do not have admin rights on the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Create namespace...")
	},
}

func init() {
	createCmd.AddCommand(createNamespaceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
