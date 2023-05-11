package delete

import (
	"github.com/spf13/cobra/doc"
	"kufast/cmd"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete root command. It cannot be executed itself but only its subcommands.
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete kufast objects.",
	Long: `The delete subcommand is a collection of all delete operations available in kufast.
Use these features to delete tenants, pods and more.`,
}

// init is a helper function from cobra to initialize the command. It sets all flags, standard values and documentation for this command.
func init() {
	cmd.RootCmd.AddCommand(deleteCmd)

}

func CreateDeleteDocs() {

	err := os.MkdirAll("./kufast.wiki/delete/", 0770)
	if err != nil {
		panic(err)
	}

	err = doc.GenMarkdownTree(deleteCmd, "./kufast.wiki/delete/")
	if err != nil {
		log.Fatal(err)
	}
}
