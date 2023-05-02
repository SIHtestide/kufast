package get

import (
	"context"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kufast/clusterOperations"
	"kufast/tools"
	"os"
)

// getCmd represents the get command
var getNamespaceCmd = &cobra.Command{
	Use:   "tenant-target <tenant-target>",
	Short: "Get information about about a namespace",
	Long:  `Get information about about a namespace.`,
	Run: func(cmd *cobra.Command, args []string) {
		//Initial config block
		clientset, _, err := tools.GetUserClient(cmd)
		if err != nil {
			tools.HandleError(err, cmd)
		}

		tenantName, err := clusterOperations.GetTenantNameFromCmd(cmd)
		if err != nil {
			tools.HandleError(err, cmd)
		}

		tenantTargetName := tenantName + "-" + args[0]

		nameSpace, err := clientset.CoreV1().Namespaces().Get(context.TODO(), tenantTargetName, metav1.GetOptions{})
		if err != nil {
			tools.HandleError(err, cmd)
		}

		quota, err := clientset.CoreV1().ResourceQuotas(tenantTargetName).Get(context.TODO(), tenantTargetName+"-limits", metav1.GetOptions{})
		if err != nil {
			tools.HandleError(err, cmd)
		}

		users, err := clientset.CoreV1().ServiceAccounts(tenantTargetName).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			tools.HandleError(err, cmd)
		}

		pods, err := clientset.CoreV1().Pods(tenantTargetName).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			tools.HandleError(err, cmd)
		}

		cpuLim, _ := quota.Spec.Hard["limits.cpu"].MarshalJSON()
		cpuReq, _ := quota.Spec.Hard["requests.cpu"].MarshalJSON()
		memLim, _ := quota.Spec.Hard["limits.memory"].MarshalJSON()
		memReq, _ := quota.Spec.Hard["requests.memory"].MarshalJSON()

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"ATTRIBUTE", "VALUE"})
		t.AppendRow(table.Row{"name", nameSpace.Name})
		t.AppendRow(table.Row{"Status", nameSpace.Status.Phase})
		t.AppendSeparator()
		t.AppendRow(table.Row{"CPU-Limit", "Limit: " + string(cpuLim) +
			"\nRequests: " + string(cpuReq)})
		t.AppendRow(table.Row{"Memory-Limit", "Limit: " + string(memLim) +
			"\nRequests: " + string(memReq)})
		t.AppendRow(table.Row{"Storage-Limit", "to be implemented"})
		t.AppendSeparator()
		t.AppendRow(table.Row{"Used CPU", quota.Status.Used.Cpu()})
		t.AppendRow(table.Row{"Used Memory", quota.Status.Used.Memory()})
		t.AppendRow(table.Row{"Available Storage", "to be implemented"})
		t.AppendSeparator()
		t.AppendRow(table.Row{"# Users", len(users.Items)})
		t.AppendRow(table.Row{"# Pods", len(pods.Items)})

		t.AppendSeparator()
		t.Render()
	},
}

func init() {
	getCmd.AddCommand(getNamespaceCmd)

	getSecretCmd.Flags().StringP("tenant", "", "", "The name of the tenant to deploy the pod to")

}
