package create

import (
	"errors"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"kufast/clusterOperations"
	"kufast/tools"
	"os"
	"time"
)

// createCmd represents the create command
var createDeploySecretCmd = &cobra.Command{
	Use:   "deploy-secret name",
	Short: "Create a new environment secret in this namespace",
	Long: `This command creates a new user and adds him to a namespace. You can select the namespace of the user.
Upon completion, the command yields the users credentials. This command will fail, if you do not have admin rights 
on the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			tools.HandleError(errors.New("too many or too few arguments provided"), cmd)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
		s.Prefix = "Creating Objects - Please wait!  "
		s.Start()

		clusterOperations.CreateDeploymentSecret(args[0], cmd)

		s.Stop()
		fmt.Println("Done!")

	},
}

func init() {
	createCmd.AddCommand(createDeploySecretCmd)

	createDeploySecretCmd.Flags().StringP("input", "", "", "Path to your .dockerfile to read credentials from.")
	createDeploySecretCmd.Flags().StringP("target", "", "", "The target for the secret (Needs to be the same as the pod using it.")
	createDeploySecretCmd.Flags().StringP("tenant", "", "", "The tenant owning the secret")
	createDeploySecretCmd.MarkFlagRequired("input")
	createDeploySecretCmd.MarkFlagFilename("input", "json")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
