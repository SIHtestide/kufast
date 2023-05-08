package list

import (
	"kufast/cmd"

	"github.com/spf13/cobra"
)

// listCmd represents the list command. It cannot be executed itself but only its subcommands.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List kufast objects",
	Long: `The list subcommand is a collection of all list operations available in kufast.
Use these features to list tenants, pods and more.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(listCmd)

}
