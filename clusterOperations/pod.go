package clusterOperations

import (
	"context"
	"errors"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kufast/objectFactory"
	"kufast/tools"
	"time"
)

func CreatePod(cmd *cobra.Command, args []string) <-chan string {
	r := make(chan string)

	go func() {
		clientset, _, err := tools.GetUserClient(cmd)
		if err != nil {
			r <- err.Error()
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

			_, err2 := clientset.CoreV1().Pods(namespaceName).Create(context.TODO(), podObject, metav1.CreateOptions{})
			if err2 != nil {
				r <- err.Error()
				return
			}

			for true {
				time.Sleep(time.Millisecond * 1000)
				pod, err := clientset.CoreV1().Pods(namespaceName).Get(context.TODO(), args[0], metav1.GetOptions{})
				if err != nil {
					r <- err.Error()
					return
				}
				if pod.Status.Phase == "Running" {
					r <- ""
					break
				}
			}
		} else {
			r <- errors.New("Invalid target for tenant").Error()
			return
		}
	}()
	return r
}
