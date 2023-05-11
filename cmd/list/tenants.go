package list

import (
	"context"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kufast/clusterOperations"
	"kufast/tools"
	"os"
)

// listTenantsCmd represents the list tenants command
var listTenantsCmd = &cobra.Command{
	Use:   "tenants",
	Short: "List all tenants in this cluster.",
	Long: `List all users in your namespace. The overview contains the name of the, the namespace where he is listed, the amount
of targets, this tenant can deploy to and the create date of this tenant.`,
	Run: func(cmd *cobra.Command, args []string) {

		s := tools.CreateStandardSpinner(tools.MESSAGE_GET_OBJECTS)

		clientset, _, err := tools.GetUserClient(cmd)
		if err != nil {
			s.Stop()
			tools.HandleError(err, cmd)
		}

		//execute request
		users, err := clientset.CoreV1().ServiceAccounts("default").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			s.Stop()
			tools.HandleError(err, cmd)
		}

		//build table
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"NAME", "NAMESPACE", "# Tenant Targets", "Created At"})
		for _, user := range users.Items {
			if user.ObjectMeta.Labels[tools.KUFAST_TENANT_LABEL] != "" {
				targets, err := clusterOperations.ListTargetsFromString(cmd, user.ObjectMeta.Labels[tools.KUFAST_TENANT_LABEL], false)
				if err != nil {

				}
				t.AppendRow(table.Row{user.Name, user.Namespace, len(targets), user.CreationTimestamp})
			}

		}

		s.Stop()
		t.AppendSeparator()
		t.Render()
	},
}

// init is a helper function from cobra to initialize the command. It sets all flags, standard values and documentation for this command.
func init() {
	listCmd.AddCommand(listTenantsCmd)

}
