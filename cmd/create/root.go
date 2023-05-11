package create

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"kufast/cmd"
	"log"
	"os"
)

// createCmd represents the create root command. It cannot be executed itself but only its subcommands.
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create new objects",
	Long: `The create subcommand is a collection of all create operations available in kufast.
Use these features to create tenants, pods and more.`,
}

// init is a helper function from cobra to initialize the command. It sets all flags, standard values and documentation for this command.
func init() {
	cmd.RootCmd.AddCommand(createCmd)

	//Enables interactive mode for all commands in create.
	createCmd.PersistentFlags().BoolP("interactive", "i", false, "Start interactive mode for the creation of this object.")

}

func CreateCreateDocs() {

	err := os.MkdirAll("./kufast.wiki/create/", 0770)
	if err != nil {
		panic(err)
	}

	err = doc.GenMarkdownTree(createCmd, "./kufast.wiki/create/")
	if err != nil {
		log.Fatal(err)
	}
}
