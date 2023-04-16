package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kufast verb object [options]",
	Short: "A small tool for creating multi-tenant environments in Kubernetes and deploy standard Docker containers to it",
	Long: `A small tool for creating a simple multi tenant environment on bare Kubernetes environments. The
tool is especially designed for people with limited Kubernetes experience, who still want to use
a Kubernetes deployment environment for their containerized applications.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringP("kubeconfig", "k", "", "Your kubeconfig to access the cluster. If not provided, we read it from $HOME/.kube/config")
	rootCmd.PersistentFlags().StringP("namespace", "n", "", "Namespace for the operation, defaults to the namespace of the default-context in your kubeconfig.")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
