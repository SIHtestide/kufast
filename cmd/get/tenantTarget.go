package get

import (
	"context"
	"errors"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kufast/clusterOperations"
	"kufast/tools"
	"os"
)

// getTenantTargetCmd represents the tenant-target command
var getTenantTargetCmd = &cobra.Command{
	Use:   "tenant-target <tenant-target>",
	Short: "Gain information on a tenant target.",
	Long:  `Gain information on a tenant target. Lists name, status, limits, their usage, and the number of included pods`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			tools.HandleError(errors.New(tools.ERROR_WRONG_NUMBER_ARGUMENTS), cmd)
		}

		s := tools.CreateStandardSpinner(tools.MESSAGE_GET_OBJECTS)

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

		pods, err := clientset.CoreV1().Pods(tenantTargetName).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			tools.HandleError(err, cmd)
		}

		cpuLim, _ := quota.Spec.Hard["limits.cpu"].MarshalJSON()
		cpuReq, _ := quota.Spec.Hard["requests.cpu"].MarshalJSON()
		memLim, _ := quota.Spec.Hard["limits.memory"].MarshalJSON()
		memReq, _ := quota.Spec.Hard["requests.memory"].MarshalJSON()
		StorageLim, _ := quota.Spec.Hard["limits.ephemeral-storage"].MarshalJSON()
		StorageReq, _ := quota.Spec.Hard["requests.ephemeral-storage"].MarshalJSON()

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
		t.AppendRow(table.Row{"Storage-Limit", "Limit: " + string(StorageLim) +
			"\nRequests: " + string(StorageReq)})
		t.AppendSeparator()
		t.AppendRow(table.Row{"Used CPU", quota.Status.Used.Cpu()})
		t.AppendRow(table.Row{"Used Memory", quota.Status.Used.Memory()})
		t.AppendRow(table.Row{"Used Storage", quota.Status.Used.Storage()})
		t.AppendRow(table.Row{"Available Storage", "to be implemented"})
		t.AppendSeparator()
		t.AppendRow(table.Row{"# Pods", len(pods.Items)})
		t.AppendSeparator()

		s.Stop()
		t.Render()
	},
}

// init is a helper function from cobra to initialize the command. It sets all flags, standard values and documentation for this command.
func init() {
	getCmd.AddCommand(getTenantTargetCmd)

	getTenantTargetCmd.Flags().StringP("tenant", "", "", tools.DOCU_FLAG_TENANT)

}
