package create

import (
	"github.com/spf13/cobra"
	"kufast/cmd"
)

// createCmd represents the create root command. It cannot be executed itself but only its subcommands.
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create new objects",
	Long: `The create subcommand is a collection of all create operations available in kufast.
Use these features to create tenants, pods and more.`,
}

func init() {
	cmd.RootCmd.AddCommand(createCmd)

	//Enables interactive mode for all commands in create.
	createCmd.PersistentFlags().BoolP("interactive", "i", false, "Start interactive mode for the creation of this object.")
}
