package update

import (
	"context"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kufast/tools"
	"os"
	"time"
)

// listCmd represents the list command
var updatePodCmd = &cobra.Command{
	Use:   "pod <pod>",
	Short: "Update an existing pod",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {

		clientset, _, err := tools.GetUserClient(cmd)
		if err != nil {
			tools.HandleError(err, cmd)

		}

		ram, _ := cmd.Flags().GetString("limit-memory")
		cpu, _ := cmd.Flags().GetString("limit-cpu")
		node, _ := cmd.Flags().GetString("target")

		namespaceName, _ := tools.GetNamespaceFromUserConfig(cmd)

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
		s.Prefix = "Creating Objects - Please wait!  "
		s.Start()

		pod, err := clientset.CoreV1().Pods(namespaceName).Get(context.TODO(), args[0], metav1.GetOptions{})
		if err != nil {
			s.Stop()
			tools.HandleError(err, cmd)

		}

		if ram != "" {
			qty, err := resource.ParseQuantity(ram)
			if err == nil {
				pod.Spec.Containers[0].Resources.Limits["memory"] = qty
				pod.Spec.Containers[0].Resources.Requests["memory"] = qty
			}
		}
		if cpu != "" {
			qty, err := resource.ParseQuantity(cpu)
			if err == nil {
				pod.Spec.Containers[0].Resources.Limits["cpu"] = qty
				pod.Spec.Containers[0].Resources.Requests["cpu"] = qty
			}
		}
		if node != "" {
			pod.Spec.NodeName = node
		}

		_, err = clientset.CoreV1().Pods(namespaceName).Update(context.TODO(), pod, metav1.UpdateOptions{})
		if err != nil {
			s.Stop()
			tools.HandleError(err, cmd)
		}

		s.Stop()
		fmt.Println("Complete!")

	},
}

func init() {
	updateCmd.AddCommand(updatePodCmd)

	updatePodCmd.Flags().BoolP("keep-alive", "", false, "Pod will be restarted upon termination.")
	updatePodCmd.Flags().StringP("memory", "", "", "Limit the RAM usage for this namespace")
	updatePodCmd.Flags().StringP("cpu", "", "", "Limit the CPU usage for this namespace")
	updatePodCmd.Flags().StringP("target", "t", "", "The name of the node to deploy the pod")
	updatePodCmd.Flags().StringArrayP("secrets", "s", []string{}, "List of secret names to be introduced in the container as environment variables. The name equals the name of the secret")

}
