package get

import (
	"context"
	"errors"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kufast/tools"
	"os"
)

// getCmd represents the get command
var getDeploySecretCmd = &cobra.Command{
	Use:   "deploy-secret <secret>",
	Short: "Gain information about a deployed pod.",
	Long:  `Gain information about a deployed pod.`,
	Run: func(cmd *cobra.Command, args []string) {
		//Initial config block
		namespaceName, err := tools.GetNamespaceFromUserConfig(cmd)
		if err != nil {
			tools.HandleError(err, cmd)
		}

		clientset, _, err := tools.GetUserClient(cmd)
		if err != nil {
			tools.HandleError(err, cmd)
		}

		//execute request
		secret, err := clientset.CoreV1().Secrets(namespaceName).Get(context.TODO(), args[0], metav1.GetOptions{})
		if err != nil {
			tools.HandleError(err, cmd)
		}

		if secret.Data[".dockerconfigjson"] == nil {
			err := errors.New("Error: This is not a deployment secret")
			tools.HandleError(err, cmd)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"ATTRIBUTE", "VALUE"})
		t.AppendRow(table.Row{"Name", secret.Name})
		t.AppendRow(table.Row{"Namespace", secret.Namespace})
		t.AppendRow(table.Row{"Data", string(secret.Data[".dockerconfigjson"])})

		t.AppendSeparator()
		t.Render()

	},
}

func init() {
	getCmd.AddCommand(getDeploySecretCmd)
}
