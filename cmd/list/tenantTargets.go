/*
MIT License

Copyright (c) 2023 Stefan Pawlowski

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
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

// listTenantTargetsCmd represents the list tenant-targets command
var listTenantTargetsCmd = &cobra.Command{
	Use:   "tenant-targets",
	Short: "List all tenant-targets of a tenant.",
	Long:  `List all tenant-targets of a tenant. The overview contains the limit information of each tenant target.`,
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

// init is a helper function from cobra to initialize the command. It sets all flags, standard values and documentation for this command.
func init() {
	listCmd.AddCommand(listTenantTargetsCmd)
	listTenantTargetsCmd.Flags().StringP("tenant", "", "", tools.DOCU_FLAG_TENANT)

}
