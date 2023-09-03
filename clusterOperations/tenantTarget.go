/*
MIT License

Copyright (c) 2023 Stefan Pawlowski

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

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

// CreateTenantTarget creates a new tenant-target
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

		target, err := GetTargetFromTargetName(cmd, targetName, tenantName, true)
		if err != nil {
			res <- err.Error()
			return
		}

		_, err = clientset.CoreV1().Namespaces().Create(context.TODO(), objectFactory.NewNamespace(tenantName, target, cmd), metav1.CreateOptions{})
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

// DeleteTenantTarget deletes a tenant-target
func DeleteTenantTarget(targetName string, tenantName string, cmd *cobra.Command) <-chan string {
	res := make(chan string)

	go func() {
		defer close(res)

		//Configblock
		clientset, _, err := tools.GetUserClient(cmd)
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

// GetTenantTarget gets a tenant-target
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

// ListTenantTarget lists a new tenant-target
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

// GetTenantTargetNameFromCmd gets the tenant targets name from the command.
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
	} else if targetName != "" {
		tenantName = tools.GetTenantFromNamespace(namespaceName)
		namespaceName = tenantName + "-" + targetName
	} else if tenantName != "" {
		defaultTargetName, err := GetTenantDefaultTargetNameFromCmd(cmd)
		if err != nil {
			return "", err
		}
		namespaceName = tenantName + "-" + defaultTargetName
	}
	return namespaceName, nil
}
