package update

import (
	"kufast/cmd"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a resource",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(updateCmd)

}
