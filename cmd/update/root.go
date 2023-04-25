package update

import (
	"fmt"
	"kufast/cmd"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a resource",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list called")
	},
}

func init() {
	cmd.RootCmd.AddCommand(updateCmd)

}
