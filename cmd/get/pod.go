package get

import (
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"kufast/clusterOperations"
	"kufast/tools"
	"os"
)

// getPodCmd represents the get pod command
var getPodCmd = &cobra.Command{
	Use:   "pod <pod>",
	Short: "Gain information about a deployed pod.",
	Long: `Gain information about a deployed pod. Output includes name, tenant-target, status, node, limits, image,
restart policy and IP-address`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			tools.HandleError(errors.New(tools.ERROR_WRONG_NUMBER_ARGUMENTS), cmd)
		}

		s := tools.CreateStandardSpinner(tools.MESSAGE_GET_OBJECTS)

		pod, err := clusterOperations.GetPod(args[0], cmd)
		if err != nil {
			tools.HandleError(err, cmd)
		}

		podEvents, err := clusterOperations.GetPodEvents(args[0], cmd)
		if err != nil {
			tools.HandleError(err, cmd)
		}

		cpuLim, _ := pod.Spec.Containers[0].Resources.Limits["cpu"].MarshalJSON()
		cpuReq, _ := pod.Spec.Containers[0].Resources.Requests["cpu"].MarshalJSON()
		memLim, _ := pod.Spec.Containers[0].Resources.Limits["memory"].MarshalJSON()
		memReq, _ := pod.Spec.Containers[0].Resources.Requests["memory"].MarshalJSON()
		storageLim, _ := pod.Spec.Containers[0].Resources.Limits["ephemeral-storage"].MarshalJSON()
		storageReq, _ := pod.Spec.Containers[0].Resources.Requests["ephemeral-storage"].MarshalJSON()

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"ATTRIBUTE", "VALUE"})
		t.AppendRow(table.Row{"Name", pod.Name})
		t.AppendRow(table.Row{"Tenant-Target", pod.Namespace})
		t.AppendRow(table.Row{"Status", pod.Status.Phase})
		t.AppendRow(table.Row{"Deployed on", pod.Spec.NodeName})
		t.AppendSeparator()
		t.AppendRow(table.Row{"CPU-Limit", "Limit: " + string(cpuLim) +
			"\nRequests: " + string(cpuReq)})
		t.AppendSeparator()
		t.AppendRow(table.Row{"Memory-Limit", "Limit: " + string(memLim) +
			"\nRequests: " + string(memReq)})
		t.AppendSeparator()
		t.AppendRow(table.Row{"Storage-Limit", "Limit: " + string(storageLim) +
			"\nRequests: " + string(storageReq)})
		t.AppendSeparator()
		t.AppendRow(table.Row{"Attached Storage"})
		t.AppendRow(table.Row{"Deployed Image", pod.Spec.Containers[0].Image})
		t.AppendRow(table.Row{"Restart Policy", pod.Spec.RestartPolicy})
		t.AppendRow(table.Row{"IP Address", pod.Status.PodIP})

		t.AppendSeparator()
		s.Stop()
		t.Render()

		fmt.Println("\n" + "Event Messages:")
		for _, podEvent := range podEvents {
			fmt.Println("\n" + podEvent.CreationTimestamp.String() + ":")
			fmt.Println(podEvent.Message)
			fmt.Println("Reason " + podEvent.Reason)
		}

	},
}

// init is a helper function from cobra to initialize the command. It sets all flags, standard values and documentation for this command.
func init() {
	getCmd.AddCommand(getPodCmd)

	getPodCmd.Flags().StringP("target", "", "", tools.DOCU_FLAG_TARGET)
	getPodCmd.Flags().StringP("tenant", "", "", tools.DOCU_FLAG_TENANT)
}
