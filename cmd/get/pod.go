package get

import (
	"context"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kufast/tools"
	"os"
)

// getCmd represents the get command
var getPodCmd = &cobra.Command{
	Use:   "pod <pod>",
	Short: "Gain information about a deployed pod.",
	Long:  `Gain information about a deployed pod.`,
	Run: func(cmd *cobra.Command, args []string) {
		//Initial config block
		ns, err := tools.GetNamespaceFromUserConfig(cmd)
		if err != nil {
			tools.HandleError(err, cmd)
		}

		clientset, _, err := tools.GetUserClient(cmd)
		if err != nil {
			tools.HandleError(err, cmd)
		}

		//execute request
		pod, err := clientset.CoreV1().Pods(ns).Get(context.TODO(), args[0], metav1.GetOptions{})
		if err != nil {
			tools.HandleError(err, cmd)
		}

		cpuLim, _ := pod.Spec.Containers[0].Resources.Limits["cpu"].MarshalJSON()
		cpuReq, _ := pod.Spec.Containers[0].Resources.Requests["cpu"].MarshalJSON()
		memLim, _ := pod.Spec.Containers[0].Resources.Limits["memory"].MarshalJSON()
		memReq, _ := pod.Spec.Containers[0].Resources.Requests["memory"].MarshalJSON()

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"ATTRIBUTE", "VALUE"})
		t.AppendRow(table.Row{"Name", pod.Name})
		t.AppendRow(table.Row{"Namespace", pod.Namespace})
		t.AppendRow(table.Row{"Status", pod.Status.Phase})
		t.AppendSeparator()
		t.AppendRow(table.Row{"CPU-Limit", "Limit: " + string(cpuLim) +
			"\nRequests: " + string(cpuReq)})
		t.AppendSeparator()
		t.AppendRow(table.Row{"Memory-Limit", "Limit: " + string(memLim) +
			"\nRequests: " + string(memReq)})
		t.AppendSeparator()
		t.AppendRow(table.Row{"Attached Storage", "None / to be implemented"})
		t.AppendRow(table.Row{"Deployed Image", pod.Spec.Containers[0].Image})
		t.AppendRow(table.Row{"Restart Policy", pod.Spec.RestartPolicy})

		t.AppendSeparator()
		t.Render()

	},
}

func init() {
	getCmd.AddCommand(getPodCmd)
}
