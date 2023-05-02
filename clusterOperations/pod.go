package clusterOperations

import (
	"context"
	"errors"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kufast/objectFactory"
	"kufast/tools"
	"time"
)

func CreatePod(cmd *cobra.Command, args []string) <-chan string {
	res := make(chan string)

	go func() {
		defer close(res)
		clientset, _, err := tools.GetUserClient(cmd)
		if err != nil {
			res <- err.Error()
			return
		}

		ram, _ := cmd.Flags().GetString("memory")
		cpu, _ := cmd.Flags().GetString("cpu")
		storage, _ := cmd.Flags().GetString("storage")
		keepAlive, _ := cmd.Flags().GetBool("keep-alive")
		target, _ := cmd.Flags().GetString("target")
		secrets, _ := cmd.Flags().GetStringArray("secrets")
		deploySecret, _ := cmd.Flags().GetString("deploy-secret")

		namespaceName, err := GetTenantTargetNameFromCmd(cmd)

		if target == "" || IsValidTarget(cmd, target, false) {

			podObject := objectFactory.NewPod(args[0], args[1], namespaceName, secrets, deploySecret, cpu, ram, storage, keepAlive)

			_, err := clientset.CoreV1().Pods(namespaceName).Create(context.TODO(), podObject, metav1.CreateOptions{})
			if err != nil {
				res <- err.Error()
				return
			}

			for true {
				time.Sleep(time.Millisecond * 1000)
				pod, err := clientset.CoreV1().Pods(namespaceName).Get(context.TODO(), args[0], metav1.GetOptions{})
				if err != nil {
					res <- err.Error()
					return
				}
				if pod.Status.Phase == "Running" {
					res <- ""
					break
				}
			}
		} else {
			res <- errors.New("Invalid target for tenant").Error()
			return
		}
	}()
	return res
}

func DeletePod(cmd *cobra.Command, pod string) <-chan string {
	res := make(chan string)

	go func() {
		defer close(res)
		clientset, _, err := tools.GetUserClient(cmd)
		if err != nil {
			res <- err.Error()
			return
		}

		namespaceName, err := GetTenantTargetNameFromCmd(cmd)
		if err != nil {
			res <- err.Error()
			return
		}

		err = clientset.CoreV1().Pods(namespaceName).Delete(context.TODO(), pod, metav1.DeleteOptions{})
		if err != nil {
			res <- err.Error()
			return
		}

		//Check for the pod been deleted from the system
		for true {
			time.Sleep(time.Millisecond * 250)
			_, err := clientset.CoreV1().Pods(namespaceName).Get(context.TODO(), pod, metav1.GetOptions{})
			if err != nil {
				res <- ""
				break
			}
		}

	}()

	return res
}

func GetPod(podName string, cmd *cobra.Command) (*v1.Pod, error) {
	//Initial config block
	namespaceName, err := GetTenantTargetNameFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	clientset, _, err := tools.GetUserClient(cmd)
	if err != nil {
		return nil, err
	}

	//execute request
	pod, err := clientset.CoreV1().Pods(namespaceName).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return pod, nil
}

func ListTenantPods(cmd *cobra.Command) ([]v1.Pod, error) {

	clientset, _, err := tools.GetUserClient(cmd)
	if err != nil {
		return nil, err
	}

	targets, err := ListTargetsFromCmd(cmd, false)
	if err != nil {
		return nil, err
	}

	tenantName, err := GetTenantNameFromCmd(cmd)
	if err != nil {
		return nil, err
	}
	var results []v1.Pod
	for _, target := range targets {
		list, err := clientset.CoreV1().Pods(tenantName+"-"+target.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		results = append(results, list.Items...)
	}

	return results, nil

}
