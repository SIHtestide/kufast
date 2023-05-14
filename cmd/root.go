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
