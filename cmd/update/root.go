package update

import (
	"kufast/cmd"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update kufast objects",
	Long: `The update subcommand is a collection of all update operations available in kufast.
Use these features to update tenants and more.`,
}

func init() {
	cmd.RootCmd.AddCommand(updateCmd)

}
