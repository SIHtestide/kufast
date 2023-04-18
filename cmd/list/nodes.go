package list

import (
	"context"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kufast/tools"
	"os"
	"strings"
)

// listCmd represents the list command
var listNodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "List all nodes in the cluster.",
	Long:  `List all nodes in the cluster. The overview contains information about the status of each node.`,
	Run: func(cmd *cobra.Command, args []string) {
		clientset, _, err := tools.GetUserClient(cmd)
		if err != nil {
			tools.HandleError(err, cmd)
		}

		//execute request
		nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			tools.HandleError(err, cmd)
		}

		//build table
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"NAME", "STATUS", "TOTAL CPU", "TOTAL RAM"})
		for _, node := range nodes.Items {
			strCond := ""
			for _, condition := range node.Status.Conditions {
				if condition.Status != "False" {
					strCond = strCond + string(condition.Type) + ", "
				}
			}
			strCond = strings.TrimSuffix(strCond, ", ")
			t.AppendRow(table.Row{node.Name, strCond, node.Status.Capacity.Cpu().String(), node.Status.Capacity.Memory().String()})

		}
		t.AppendSeparator()
		t.Render()
	},
}

func init() {
	listCmd.AddCommand(listNodesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
