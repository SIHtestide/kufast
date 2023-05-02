package list

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"kufast/clusterOperations"
	"kufast/tools"
	"os"
)

// listCmd represents the list command
var listPodsCmd = &cobra.Command{
	Use:   "pods",
	Short: "List all pods in a namespace",
	Long: `List all pods in your namespace. The overview contains information about the status of your pod and its name
To gain further information see the kubectl get pod command.`,
	Run: func(cmd *cobra.Command, args []string) {

		s := tools.CreateStandardSpinner(tools.MESSAGE_GET_OBJECTS)

		pods, err := clusterOperations.ListTenantPods(cmd)
		if err != nil {
			s.Stop()
			tools.HandleError(err, cmd)
		}

		//build table
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"NAME", "NAMESPACE", "STATUS", "MESSAGE"})
		for _, pod := range pods {
			t.AppendRow(table.Row{pod.Name, pod.Namespace, pod.Status.Phase, pod.Status.Message})

		}
		s.Stop()
		t.AppendSeparator()
		t.Render()
	},
}

func init() {
	listCmd.AddCommand(listPodsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
