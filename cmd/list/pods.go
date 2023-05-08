package list

import (
	"context"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		t.AppendHeader(table.Row{"NAME", "NAMESPACE", "STATUS", "POD MESSAGE", "CONTAINER MESSAGE"})
		for _, pod := range pods {
			t.AppendRow(table.Row{pod.Name, pod.Namespace, pod.Status.Phase, pod.Status.Message, ""})

		}
		s.Stop()
		t.AppendSeparator()
		t.Render()

		clientset, _, _ := tools.GetUserClient(cmd)

		events, _ := clientset.CoreV1().Events("testtenant3-w2").List(context.TODO(), metav1.ListOptions{})
		for _, event := range events.Items {
			fmt.Println(event)
		}
	},
}

func init() {
	listCmd.AddCommand(listPodsCmd)
	listPodsCmd.Flags().StringP("tenant", "", "", "The name of the tenant to deploy the pod to")
}
