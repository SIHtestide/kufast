package list

import (
	"context"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kufast/tools"
	"os"
)

// listCmd represents the list command
var listnamespacesCmd = &cobra.Command{
	Use:   "namespaces",
	Short: "List all namespaces in the cluster.",
	Long:  `List all nodes in the cluster. The overview contains information about the status of each node.`,
	Run: func(cmd *cobra.Command, args []string) {
		clientset, _, err := tools.GetUserClient(cmd)
		if err != nil {
			tools.HandleError(err, cmd)
		}

		//execute request
		namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			tools.HandleError(err, cmd)
		}

		//build table
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"NAME", "STATUS", "CPU Limit", "Memory Limit", "Storage Limit"})
		for _, namespace := range namespaces.Items {

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
		t.AppendSeparator()
		t.Render()
	},
}

func init() {
	listCmd.AddCommand(listnamespacesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
