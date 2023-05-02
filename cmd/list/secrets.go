package list

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"kufast/clusterOperations"
	"kufast/tools"
	"os"
)

// listCmd represents the list command
var listSecretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "List all secrets in a namespace",
	Long: `List all secrets in your namespace. The overview contains information about the nameof the secret and its creation date
To gain further information see the kubectl get pod command.`,
	Run: func(cmd *cobra.Command, args []string) {

		s := tools.CreateStandardSpinner(tools.MESSAGE_GET_OBJECTS)

		secrets, err := clusterOperations.ListSecrets(cmd)
		if err != nil {
			s.Stop()
			tools.HandleError(err, cmd)
		}

		//build table
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"NAME", "NAMESPACE", "TYPE", "CREATED AT"})
		for _, secret := range secrets {
			t.AppendRow(table.Row{secret.Name, secret.Namespace, secret.Type, secret.CreationTimestamp})
		}

		s.Stop()
		t.AppendSeparator()
		t.Render()
	},
}

func init() {
	listCmd.AddCommand(listSecretsCmd)

	listSecretsCmd.Flags().StringP("tenant", "", "", "The name of the tenant to deploy the pod to")

}
