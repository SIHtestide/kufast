package get

import (
	"github.com/spf13/cobra/doc"
	"kufast/cmd"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// getCmd represents the get root command. It cannot be executed itself but only its subcommands.
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get existing kufast objects.",
	Long: `The get subcommand is a collection of all get operations available in kufast.
Use these features to get tenants, pods and more.`,
}

// init is a helper function from cobra to initialize the command. It sets all flags, standard values and documentation for this command.
func init() {
	cmd.RootCmd.AddCommand(getCmd)

}

func CreateGetDocs() {
	err := os.MkdirAll("./kufast.wiki/get/", 0770)
	if err != nil {
		panic(err)
	}

	err = doc.GenMarkdownTree(getCmd, "./kufast.wiki/get/")
	if err != nil {
		log.Fatal(err)
	}
}
