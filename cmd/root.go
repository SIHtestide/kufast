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
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"os"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "kufast verb object [options]",
	Short: "A small tool for creating multi-tenant environments in Kubernetes and deploy standard Docker containers to it",
	Long: `A small tool for creating a simple multi tenant environment on bare Kubernetes environments. The
tool is especially designed for people with limited Kubernetes experience, who still want to use
a Kubernetes deployment environment for their containerized applications.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// init is a helper function from cobra to initialize the command. It sets all flags, standard values and documentation for this command.
func init() {
	RootCmd.PersistentFlags().StringP("kubeconfig", "k", "", "Your kubeconfig to access the cluster. If not provided, we read it from $HOME/.kube/config")

}

func CreateRootDocs(linkH func(string) string) {
	os.MkdirAll("./kufast.wiki", 0770)
	out, err := os.Create("./kufast.wiki/root.md")
	if err != nil {
		return
	}

	defer func() {
		err := out.Close()
		if err != nil {
			panic(err)
		}
	}()

	err = doc.GenMarkdownCustom(RootCmd, out, linkH)
	if err != nil {
		panic(err)
	}

}
