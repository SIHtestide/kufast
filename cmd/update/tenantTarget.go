package update

import (
	"context"
	"errors"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kufast/clusterOperations"
	"kufast/objectFactory"
	"kufast/tools"
	"os"
	"time"
)

// updateTenantTargetCmd represents the update tenant-target command
var updateTenantTargetCmd = &cobra.Command{
	Use:   "tenant-target <tenant target>",
	Short: "Update memory, CPU and storage capabilities of a tenant target.",
	Long: "Update memory, CPU and storage capabilities of a tenant target. " +
		"Also updates the role scheme to the latest version of kufast.",
	Run: func(cmd *cobra.Command, args []string) {

		//Check that exactly one arg has been provided (the namespace)
		if len(args) < 1 {
			tools.HandleError(errors.New(tools.ERROR_WRONG_NUMBER_ARGUMENTS), cmd)
		}

		//Initial config block
		clientset, _, err := tools.GetUserClient(cmd)
		if err != nil {
			tools.HandleError(err, cmd)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
		s.Prefix = tools.MESSAGE_UPDATE_OBJECTS
		s.Start()

		tenantName, err := clusterOperations.GetTenantNameFromCmd(cmd)
		if err != nil {
			s.Stop()
			tools.HandleError(err, cmd)
		}

		tenantTargetName := tenantName + "-" + args[0]

		//Get Current Namespace
		namespace, err := clientset.CoreV1().Namespaces().Get(context.TODO(), tenantTargetName, metav1.GetOptions{})
		if err != nil {
			s.Stop()
			tools.HandleError(err, cmd)
		}

		//Get quotas for namespace
		quota, err := clientset.CoreV1().ResourceQuotas(tenantTargetName).Get(context.TODO(), tenantTargetName+"-limits", metav1.GetOptions{})
		if err != nil {
			s.Stop()
			tools.HandleError(err, cmd)
		}

		//Get Networkpolicy for namespace
		nps, err := clientset.NetworkingV1().NetworkPolicies(tenantTargetName).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			s.Stop()
			tools.HandleError(err, cmd)
		}

		ram, _ := cmd.Flags().GetString("memory")
		cpu, _ := cmd.Flags().GetString("cpu")
		target, _ := cmd.Flags().GetString("target")

		if namespace.ObjectMeta.Annotations == nil {
			//No annotations have been provided, need to create them
			annotations := map[string]string{}
			namespace.ObjectMeta.Annotations = annotations
		}
		namespace.ObjectMeta.Annotations["scheduler.alpha.kubernetes.io/node-selector"] = "target=" + target

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

		if len(nps.Items) == 0 {
			//Keine Policy, erstelle Policy
			_, _ = clientset.NetworkingV1().NetworkPolicies(args[0]).Create(context.TODO(), objectFactory.NewNetworkPolicy(args[0], tenantName), metav1.CreateOptions{})
		} else if len(nps.Items) == 1 {
			//Update Policy
			_, _ = clientset.NetworkingV1().NetworkPolicies(args[0]).Update(context.TODO(), objectFactory.NewNetworkPolicy(args[0], tenantName), metav1.UpdateOptions{})
		} else {
			s.Stop()
			fmt.Println("More than one Network policy detected! ignoring..")
			s.Start()
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

		_, err = clientset.CoreV1().Namespaces().Update(context.TODO(), namespace, metav1.UpdateOptions{})
		if err != nil {
			s.Stop()
			tools.HandleError(err, cmd)
		}

		s.Stop()
		fmt.Println("Complete!")

	},
}

// init is a helper function from cobra to initialize the command. It sets all flags, standard values and documentation for this command.
func init() {
	updateCmd.AddCommand(updateTenantTargetCmd)

	updateTenantTargetCmd.Flags().StringP("memory", "", "4Gi", "Limit the RAM usage for this namespace")
	updateTenantTargetCmd.Flags().StringP("cpu", "", "2", "Limit the CPU usage for this namespace")
	updateTenantTargetCmd.Flags().StringP("tenant", "t", "", "The tenant owning this namespace. Matching tenants will be connected.")
	_ = updateTenantTargetCmd.MarkFlagRequired("tenant")

}
