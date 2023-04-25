package update

import (
	"context"
	"errors"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kufast/objectFactory"
	"kufast/tools"
	"os"
	"time"
)

// listCmd represents the list command
var updateNamespaceCmd = &cobra.Command{
	Use:   "namespace <namespace>",
	Short: "Update Memory and CPU capabilities. Updates the role scheme to the latest version.",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {

		//Check that exactly one arg has been provided (the namespace)
		if len(args) < 1 {
			tools.HandleError(errors.New("Too few arguments provided."), cmd)
		}

		//Initial config block
		clientset, _, err := tools.GetUserClient(cmd)
		if err != nil {
			tools.HandleError(err, cmd)
		}

		//Get Current Namespace to verify it is present
		_, err = clientset.CoreV1().Namespaces().Get(context.TODO(), args[0], metav1.GetOptions{})
		if err != nil {
			tools.HandleError(err, cmd)
		}

		//Get quotas for namespace
		quota, err := clientset.CoreV1().ResourceQuotas(args[0]).Get(context.TODO(), args[0]+"-limits", metav1.GetOptions{})
		if err != nil {
			tools.HandleError(err, cmd)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
		s.Prefix = "Creating Objects - Please wait!  "
		s.Start()

		ram, _ := cmd.Flags().GetString("memory")
		cpu, _ := cmd.Flags().GetString("cpu")

		if ram != "" {
			qty, err := resource.ParseQuantity(ram)
			if err == nil {
				quota.Spec.Hard["limits.memory"] = qty
				quota.Spec.Hard["requests.memory"] = qty
			}
		}
		if cpu != "" {
			qty, err := resource.ParseQuantity(cpu)
			if err == nil {
				quota.Spec.Hard["limits.cpu"] = qty
				quota.Spec.Hard["requests.cpu"] = qty
			}
		}

		//Create current role scheme to update namespace
		role := objectFactory.NewRole(args[0])

		//Apply changes
		_, err = clientset.CoreV1().ResourceQuotas(args[0]).Update(context.TODO(), quota, metav1.UpdateOptions{})
		if err != nil {
			s.Stop()
			tools.HandleError(err, cmd)
		}
		_, err = clientset.RbacV1().Roles(args[0]).Update(context.TODO(), role, metav1.UpdateOptions{})
		if err != nil {
			s.Stop()
			tools.HandleError(err, cmd)
		}
		s.Stop()
		fmt.Println("Complete!")

	},
}

func init() {
	updateCmd.AddCommand(updateNamespaceCmd)

	updateNamespaceCmd.Flags().StringP("memory", "", "4Gi", "Limit the RAM usage for this namespace")
	updateNamespaceCmd.Flags().StringP("cpu", "", "2", "Limit the CPU usage for this namespace")

}
