package clusterOperations

import (
	"context"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kufast/objectFactory"
	"kufast/tools"
	"time"
)

func CreateTenantTarget(tenantName string, targetName string, cmd *cobra.Command) <-chan string {
	res := make(chan string)

	go func() {
		defer close(res)

		clientset, _, err := tools.GetUserClient(cmd)
		if err != nil {
			res <- err.Error()
			return
		}

		newNamespaceName := tenantName + "-" + targetName

		ram, _ := cmd.Flags().GetString("memory")
		cpu, _ := cmd.Flags().GetString("cpu")
		storage, _ := cmd.Flags().GetString("storage")
		minStorage, _ := cmd.Flags().GetString("storage-min")
		pods, _ := cmd.Flags().GetString("pods")

		tenant, err := GetTenantFromString(cmd, tenantName)
		if err != nil {
			res <- err.Error()
			return
		}

		err = AddTargetToTenant(cmd, targetName, tenant)
		if err != nil {
			res <- err.Error()
			return
		}

		_, err = clientset.CoreV1().Namespaces().Create(context.TODO(), objectFactory.NewNamespace(tenantName, targetName, cmd), metav1.CreateOptions{})
		if err != nil {
			res <- err.Error()
			return
		}

		for true {
			newNamespace, err := clientset.CoreV1().Namespaces().Get(context.TODO(), tenantName+"-"+targetName, metav1.GetOptions{})
			if err != nil {
				res <- err.Error()
				return
			}
			if newNamespace.Status.Phase == "Active" {
				break
			}
			time.Sleep(time.Millisecond * 250)
		}

		_, err = clientset.CoreV1().ResourceQuotas(newNamespaceName).Create(context.TODO(), objectFactory.NewResourceQuota(newNamespaceName, ram, cpu, storage, pods), metav1.CreateOptions{})
		if err != nil {
			res <- err.Error()
			return
		}

		_, err = clientset.RbacV1().Roles(newNamespaceName).Create(context.TODO(), objectFactory.NewRole(newNamespaceName), metav1.CreateOptions{})
		if err != nil {
			res <- err.Error()
			return
		}

		_, err = clientset.CoreV1().LimitRanges(newNamespaceName).Create(context.TODO(), objectFactory.NewLimitRange(newNamespaceName, minStorage, storage), metav1.CreateOptions{})
		if err != nil {
			res <- err.Error()
			return
		}

		_, err = clientset.NetworkingV1().NetworkPolicies(newNamespaceName).Create(context.TODO(), objectFactory.NewNetworkPolicy(newNamespaceName, tenantName), metav1.CreateOptions{})
		if err != nil {
			res <- err.Error()
			return
		}

		_, err = clientset.RbacV1().RoleBindings(newNamespaceName).Create(context.TODO(), objectFactory.NewTenantRolebinding(newNamespaceName, tenantName), metav1.CreateOptions{})
		if err != nil {
			res <- err.Error()
			return
		}

		res <- ""
	}()
	return res

}

func DeleteTenantTarget(targetName string, cmd *cobra.Command) <-chan string {
	res := make(chan string)

	go func() {
		defer close(res)

		//Configblock
		clientset, _, err := tools.GetUserClient(cmd)
		if err != nil {
			res <- err.Error()
			return
		}

		tenantName, err := cmd.Flags().GetString("tenant")
		if err != nil {
			res <- err.Error()
			return
		}

		err = DeleteTargetFromTenant(targetName, cmd)
		if err != nil {
			res <- err.Error()
			return
		}

		err = clientset.CoreV1().Namespaces().Delete(context.TODO(), tenantName+"-"+targetName, metav1.DeleteOptions{})
		if err != nil {
			res <- err.Error()
			return
		}

		res <- ""

	}()
	return res

}

func GetTenantTarget(tenantName string, targetName string, cmd *cobra.Command) (*v1.Namespace, error) {

	//Configblock
	clientset, _, err := tools.GetUserClient(cmd)
	if err != nil {
		return nil, err
	}

	tenantTarget, err := clientset.CoreV1().Namespaces().Get(context.TODO(), tenantName+"-"+targetName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return tenantTarget, nil

}

func ListTenantTargets(tenantName string, cmd *cobra.Command) ([]*v1.Namespace, error) {

	tenantTargets, err := ListTargetsFromString(cmd, tenantName, false)
	if err != nil {
		return nil, err
	}

	var tenantTargetObjects []*v1.Namespace

	for _, target := range tenantTargets {
		tenantTarget, err := GetTenantTarget(tenantName, target.Name, cmd)
		if err != nil {
			return nil, err
		}
		tenantTargetObjects = append(tenantTargetObjects, tenantTarget)
	}

	return tenantTargetObjects, nil

}

func UpdateTenantTarget(target string, cmd *cobra.Command) (*v1.Namespace, error) {

}

func GetTenantTargetNameFromCmd(cmd *cobra.Command) (string, error) {

	namespaceName, err := tools.GetNamespaceFromUserConfig(cmd)
	if err != nil {
		return "", err
	}
	tenantName, err := cmd.Flags().GetString("tenant")
	if err != nil {
		return "", err
	}
	targetName, err := cmd.Flags().GetString("target")
	if err != nil {
		return "", err
	}

	if tenantName != "" && targetName != "" {
		namespaceName = tenantName + "-" + targetName
	} else if tenantName != "" {
		tenantName = tools.GetTenantFromNamespace(namespaceName)
		namespaceName = tenantName + "-" + targetName
	} else if targetName != "" {
		defaultTargetName, err := GetTenantDefaultTargetNameFromCmd(cmd)
		if err != nil {
			return "", err
		}
		namespaceName = tenantName + "-" + defaultTargetName
	}
	return namespaceName, nil
}
