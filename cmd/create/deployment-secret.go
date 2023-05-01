package create

import (
	"context"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kufast/objectFactory"
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
		clientset, _, err := tools.GetUserClient(cmd)
		if err != nil {
			tools.HandleError(err, cmd)
		}

		namespaceName := tools.GetDeploymentNamespace(cmd)
		fileName, _ := cmd.Flags().GetString("input")

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
		s.Prefix = "Creating Objects - Please wait!  "
		s.Start()

		creds, err := os.ReadFile(fileName)
		if err != nil {
			s.Stop()
			tools.HandleError(err, cmd)
		}

		deploymentSecretObject := objectFactory.NewDeploymentSecret(namespaceName, args[0], creds)

		_, err = clientset.CoreV1().Secrets(namespaceName).Create(context.TODO(), deploymentSecretObject, v1.CreateOptions{})
		if err != nil {
			s.Stop()
			tools.HandleError(err, cmd)
		}

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
