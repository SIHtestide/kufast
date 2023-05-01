package tools

import (
	"context"
	"errors"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/strings/slices"
	"strings"
)

type Target struct {
	Name       string
	AccessType string
}

func IsValidTarget(cmd *cobra.Command, target string, all bool) bool {
	if strings.Contains(target, "_") {
		return false
	}

	targets := ListTargets(cmd, all)
	for _, t := range targets {
		if t.Name == target {
			return true
		}
	}
	return false
}

func GetTenantDefaultTargetName(tenant string, cmd *cobra.Command) string {
	clientset, _, _ := GetUserClient(cmd)

	user, _ := clientset.CoreV1().ServiceAccounts("default").Get(context.TODO(), tenant+"-user", metav1.GetOptions{})
	return user.ObjectMeta.Labels["kufast/defaultTarget"]
}
func GetTargetFromTargetName(cmd *cobra.Command, target string, all bool) Target {
	targets := ListTargets(cmd, all)
	for _, t := range targets {
		if t.Name == target {
			return t
		}
	}
	return Target{}
}

func AddTargetToTenant(cmd *cobra.Command, targetName string, user *v1.ServiceAccount) *v1.ServiceAccount {
	if IsValidTarget(cmd, targetName, true) {
		target := GetTargetFromTargetName(cmd, targetName, true)
		if target.AccessType == "node" {
			user.ObjectMeta.Labels["kufast.nodeAccess/"+targetName] = "true"
		} else {
			user.ObjectMeta.Labels["kufast.groupAccess/"+targetName] = "true"
		}
	}

	// Populate default label if possible
	if user.ObjectMeta.Labels["kufast/defaultTarget"] == "" {
		user.ObjectMeta.Labels["kufast/defaultTarget"] = targetName
	}
	return user
}

func ListTargets(cmd *cobra.Command, all bool) []Target {

	clientset, _, _ := GetUserClient(cmd)
	namespaceName, _ := GetNamespaceFromUserConfig(cmd)
	var results []Target

	//Do we want the target of the user or all?
	if all {
		//This information is only available by parsing the nodes
		nodes, _ := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		var groups []string
		for _, node := range nodes.Items {
			//Append node target
			results = append(results, Target{
				Name:       node.ObjectMeta.Labels["kubernetes.io/hostname"],
				AccessType: "node",
			})
			for key, elem := range node.ObjectMeta.Labels {
				if strings.Contains("kufast.group/", key) && elem != "false" && !slices.Contains(groups, strings.TrimPrefix(key, "kufast.group/")) {
					groups = append(groups, strings.TrimPrefix(key, "kufast.group/"))
				}
			}
		}
		for _, target := range groups {
			if target != "" {
				results = append(results, Target{
					Name:       target,
					AccessType: "group",
				})
			}
		}

	} else {
		//Get the information from the tenant
		tenant, _ := cmd.Flags().GetString("tenant")
		if tenant == "" {
			tenant = GetTenantFromNamespace(namespaceName)
		}

		user, _ := clientset.CoreV1().ServiceAccounts("default").Get(context.TODO(), tenant+"-user", metav1.GetOptions{})

		for key, elem := range user.ObjectMeta.Labels {
			if strings.Contains("kufast.groupAccess/", key) && elem != "false" {
				results = append(results, Target{
					Name:       strings.TrimPrefix(key, "kufast.groupAccess/"),
					AccessType: "group",
				})
			} else if strings.Contains("kufast.nodeAccess/", key) && elem != "false" {
				results = append(results, Target{
					Name:       strings.TrimPrefix(key, "kufast.nodeAccess/"),
					AccessType: "node",
				})
			}
		}
	}
	return results

}

func SetTargetGroupToNodes(targetName string, targetNodes []string, cmd *cobra.Command) error {
	clientset, _, err := GetUserClient(cmd)
	if err != nil {
		return errors.New(err.Error())
	}

	nodeList, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return errors.New(err.Error())
	}

	if !IsValidTarget(cmd, targetName, true) {
		for _, node := range nodeList.Items {
			if slices.Contains(targetNodes, node.Name) {
				node.Labels["kufast.group/"+targetName] = "true"
			} else {
				node.Labels["kufast.group/"+targetName] = "false"
			}
			_, err = clientset.CoreV1().Nodes().Update(context.TODO(), &node, metav1.UpdateOptions{})
			if err != nil {
				return errors.New(err.Error())
			}
		}
	}

	return nil
}
