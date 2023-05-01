package clusterOperations

import (
	"context"
	"errors"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/strings/slices"
	"kufast/tools"
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

	targets, err := ListTargets(cmd, all)
	if err != nil {
		return false
	}
	for _, t := range targets {
		if t.Name == target {
			return true
		}
	}
	return false
}

func GetTargetFromTargetName(cmd *cobra.Command, target string, all bool) (Target, error) {
	targets, err := ListTargets(cmd, all)
	if err != nil {
		return Target{}, err
	}
	for _, t := range targets {
		if t.Name == target {
			return t, nil
		}
	}
	return Target{}, errors.New("the target does not exist or the tenant has no access to the target")
}

func ListTargets(cmd *cobra.Command, all bool) ([]Target, error) {

	clientset, _, err := tools.GetUserClient(cmd)
	if err != nil {
		return nil, err
	}
	var results []Target

	//Do we want the target of the user or all?
	if all {
		//This information is only available by parsing the nodes
		nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		var groups []string
		for _, node := range nodes.Items {
			//Append node target
			results = append(results, Target{
				Name:       node.ObjectMeta.Labels[tools.KUFAST_NODE_HOSTNAME_LABEL],
				AccessType: "node",
			})
			for key, elem := range node.ObjectMeta.Labels {
				if strings.Contains(tools.KUFAST_NODE_GROUP_LABEL, key) && elem != "false" && !slices.Contains(groups, strings.TrimPrefix(key, tools.KUFAST_NODE_GROUP_LABEL)) {
					groups = append(groups, strings.TrimPrefix(key, tools.KUFAST_NODE_GROUP_LABEL))
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
		namespaceName, _ := tools.GetNamespaceFromUserConfig(cmd)
		tenant, _ := cmd.Flags().GetString("tenant")
		if tenant == "" {
			tenant = tools.GetTenantFromNamespace(namespaceName)
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
	return results, nil

}

func SetTargetGroupToNodes(targetName string, targetNodes []string, cmd *cobra.Command) error {
	clientset, _, err := tools.GetUserClient(cmd)
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
				node.ObjectMeta.Labels["kufast.group/"+targetName] = "true"
			} else {
				node.ObjectMeta.Labels["kufast.group/"+targetName] = "false"
			}
			_, err = clientset.CoreV1().Nodes().Update(context.TODO(), &node, metav1.UpdateOptions{})
			if err != nil {
				return errors.New(err.Error())
			}
		}
	}

	return nil
}

func DeleteTargetGroupFromNodes(targetName string, cmd *cobra.Command) error {
	clientset, _, err := tools.GetUserClient(cmd)
	if err != nil {
		return errors.New(err.Error())
	}

	nodeList, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return errors.New(err.Error())
	}
	if IsValidTarget(cmd, targetName, true) {
		for _, node := range nodeList.Items {
			delete(node.ObjectMeta.Labels, "kufast.group/"+targetName)
			_, err = clientset.CoreV1().Nodes().Update(context.TODO(), &node, metav1.UpdateOptions{})
			if err != nil {
				return errors.New(err.Error())
			}
		}
	}
	return nil
}
