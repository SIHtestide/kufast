package create

import (
	"context"
	b64 "encoding/base64"
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
	Use:   "secret name",
	Short: "Create a new environment secret in this namespace",
	Long: `This command creates a new user and adds him to a namespace. You can select the namespace of the user.
Upon completion, the command yields the users credentials. This command will fail, if you do not have admin rights 
on the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		clientset, _, err := tools.GetUserClient(cmd)
		if err != nil {
			tools.HandleError(err, cmd)
		}

		namespaceName, _ := tools.GetNamespaceFromUserConfig(cmd)
		fileName, _ := cmd.Flags().GetString("input")

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
		s.Prefix = "Deleting Objects - Please wait!  "
		s.Start()

		buf, err := os.ReadFile(fileName)
		if err != nil {
			s.Stop()
			tools.HandleError(err, cmd)
		}

		credentials := b64.StdEncoding.EncodeToString(buf)
		deploymentSecretObject := objectFactory.NewDeploymentSecret(namespaceName, args[0], credentials)

		clientset.CoreV1().Secrets(namespaceName).Create(context.TODO(), deploymentSecretObject, v1.CreateOptions{})

		s.Stop()
		fmt.Println("Complete!")

	},
}

func init() {
	createCmd.AddCommand(createDeploySecretCmd)

	createDeploySecretCmd.Flags().StringP("input", "i", "", "Path to your .dockerfile to read credentials from.")
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
