package delete

import (
	"kufast/cmd"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete root command. It cannot be executed itself but only its subcommands.
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete kufast objects.",
	Long: `The delete subcommand is a collection of all delete operations available in kufast.
Use these features to delete tenants, pods and more.`,
}

func init() {
	cmd.RootCmd.AddCommand(deleteCmd)

}
