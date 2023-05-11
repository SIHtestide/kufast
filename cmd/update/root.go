package update

import (
	"github.com/spf13/cobra/doc"
	"kufast/cmd"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update kufast objects",
	Long: `The update subcommand is a collection of all update operations available in kufast.
Use these features to update tenants and more.`,
}

// init is a helper function from cobra to initialize the command. It sets all flags, standard values and documentation for this command.
func init() {
	cmd.RootCmd.AddCommand(updateCmd)

}

func CreateUpdateDocs(fileP func(string) string, linkH func(string) string) {

	err := os.MkdirAll("./kufast.wiki/update/", 0770)
	if err != nil {
		panic(err)
	}

	err = doc.GenMarkdownTreeCustom(updateCmd, "./kufast.wiki/update/", fileP, linkH)
	if err != nil {
		log.Fatal(err)
	}
}
