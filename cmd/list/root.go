/*
MIT License

Copyright (c) 2023 Stefan Pawlowski

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package list

import (
	"github.com/spf13/cobra/doc"
	"kufast/cmd"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// listCmd represents the list command. It cannot be executed itself but only its subcommands.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List kufast objects",
	Long: `The list subcommand is a collection of all list operations available in kufast.
Use these features to list tenants, pods and more.`,
}

// init is a helper function from cobra to initialize the command. It sets all flags, standard values and documentation for this command.
func init() {
	cmd.RootCmd.AddCommand(listCmd)

}

func CreateListDocs(fileP func(string) string, linkH func(string) string) {

	err := os.MkdirAll("./kufast.wiki/list/", 0770)
	if err != nil {
		panic(err)
	}

	err = doc.GenMarkdownTreeCustom(listCmd, "./kufast.wiki/list/", fileP, linkH)
	if err != nil {
		log.Fatal(err)
	}
}
