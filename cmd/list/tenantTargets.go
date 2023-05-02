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

// listCmd represents the list command
var listTenantTargetsCmd = &cobra.Command{
	Use:   "namespaces",
	Short: "List all namespaces in the cluster.",
	Long:  `List all nodes in the cluster. The overview contains information about the status of each node.`,
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

		namespaces, err := clusterOperations.ListTenantTargets(tenantName, cmd)
		if err != nil {
			tools.HandleError(err, cmd)
		}

		s := tools.CreateStandardSpinner(tools.MESSAGE_GET_OBJECTS)

		//build table
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"NAME", "STATUS", "CPU Limit", "Memory Limit", "Storage Limit"})
		for _, namespace := range namespaces {

			quotas, err := clientset.CoreV1().ResourceQuotas(namespace.Name).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				tools.HandleError(err, cmd)
			}

			var memoryQuota string
			var cpuQuota string
			var storageQuota string

			if len(quotas.Items) == 1 {
				cpuQuotaBytes, err := quotas.Items[0].Spec.Hard["limits.cpu"].MarshalJSON()
				if err != nil {
					cpuQuota = "None"
				} else {
					cpuQuota = string(cpuQuotaBytes)
				}
				memoryQuotaBytes, err := quotas.Items[0].Spec.Hard["limits.memory"].MarshalJSON()
				if err != nil {
					memoryQuota = "None"
				} else {
					memoryQuota = string(memoryQuotaBytes)
				}
				storageQuotaBytes, err := quotas.Items[0].Spec.Hard["requests.ephemeral-storage"].MarshalJSON()
				if err != nil {
					storageQuota = "None"
				} else {
					storageQuota = string(storageQuotaBytes)
				}
			} else {
				memoryQuota = "Quota missing or too many quotas"
				cpuQuota = "Quota missing or too many quotas"
			}

			t.AppendRow(table.Row{namespace.Name, namespace.Status, cpuQuota, memoryQuota, storageQuota})

		}
		s.Stop()
		t.AppendSeparator()
		t.Render()
	},
}

func init() {
	listCmd.AddCommand(listTenantTargetsCmd)
	listTenantTargetsCmd.Flags().StringP("tenant", "", "", "The name of the tenant to deploy the pod to")

}
