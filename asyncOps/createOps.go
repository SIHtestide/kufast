package asyncOps

import (
	"context"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kufast/objectFactory"
	"kufast/tools"
	"time"
)

func CreateNamespace(tenant string, target string, cmd *cobra.Command, s *spinner.Spinner) <-chan int32 {
	r := make(chan int32)

	go func() {
		defer close(r)

		clientset, _, err := tools.GetUserClient(cmd)
		if err != nil {
			tools.HandleError(err, cmd)
		}

		ram, _ := cmd.Flags().GetString("memory")
		cpu, _ := cmd.Flags().GetString("cpu")
		storage, _ := cmd.Flags().GetString("storage")
		minStorage, _ := cmd.Flags().GetString("storage-min")
		pods, _ := cmd.Flags().GetString("pods")

		_, err = clientset.CoreV1().Namespaces().Create(context.TODO(), objectFactory.NewNamespace(tenant, target, cmd), metav1.CreateOptions{})
		if err != nil {
			r <- 1
			s.Stop()
			fmt.Println(err.Error())
			s.Start()
		}

		for true {
			newNamespace, err := clientset.CoreV1().Namespaces().Get(context.TODO(), tenant+"-"+target, metav1.GetOptions{})
			if err != nil {
				r <- 1
				s.Stop()
				fmt.Println(err.Error())
				s.Start()
			}
			if newNamespace.Status.Phase == "Active" {
				break
			}
			time.Sleep(time.Millisecond * 250)
		}

		_, err = clientset.CoreV1().ResourceQuotas(tenant+"-"+target).Create(context.TODO(), objectFactory.NewResourceQuota(tenant+"-"+target, ram, cpu, storage, pods), metav1.CreateOptions{})
		if err != nil {
			r <- 1
			s.Stop()
			fmt.Println(err.Error())
			s.Start()
		}

		_, err = clientset.RbacV1().Roles(tenant+"-"+target).Create(context.TODO(), objectFactory.NewRole(tenant+"-"+target), metav1.CreateOptions{})
		if err != nil {
			r <- 1
			s.Stop()
			fmt.Println(err.Error())
			s.Start()
		}

		_, err = clientset.CoreV1().LimitRanges(tenant+"-"+target).Create(context.TODO(), objectFactory.NewLimitRange(tenant+"-"+target, minStorage, storage), metav1.CreateOptions{})
		if err != nil {
			r <- 1
			s.Stop()
			fmt.Println(err.Error())
			s.Start()
		}

		_, err = clientset.NetworkingV1().NetworkPolicies(tenant+"-"+target).Create(context.TODO(), objectFactory.NewNetworkPolicy(tenant+"-"+target, tenant), metav1.CreateOptions{})
		if err != nil {
			r <- 1
			s.Stop()
			fmt.Println(err.Error())
			s.Start()
		}

		_, err = clientset.RbacV1().RoleBindings(tenant+"-"+target).Create(context.TODO(), objectFactory.NewTenantRolebinding(tenant+"-"+target, tenant), metav1.CreateOptions{})
		if err != nil {
			r <- 1
			s.Stop()
			fmt.Println(err.Error())
			s.Start()
		}

		r <- 0
	}()
	return r

}

func CreatePod(cmd *cobra.Command, s *spinner.Spinner, args []string) <-chan int32 {
	r := make(chan int32)

	go func() {
		clientset, _, err := tools.GetUserClient(cmd)
		if err != nil {
			r <- 1
			s.Stop()
			fmt.Println(err.Error())
			s.Start()
		}

		ram, _ := cmd.Flags().GetString("memory")
		cpu, _ := cmd.Flags().GetString("cpu")
		storage, _ := cmd.Flags().GetString("storage")
		keepAlive, _ := cmd.Flags().GetBool("keep-alive")
		target, _ := cmd.Flags().GetString("target")
		secrets, _ := cmd.Flags().GetStringArray("secrets")
		deploySecret, _ := cmd.Flags().GetString("deploy-secret")

		namespaceName := tools.GetDeploymentNamespace(cmd)

		if target == "" || tools.IsValidTarget(cmd, target, false) {

			podObject := objectFactory.NewPod(args[0], args[1], namespaceName, secrets, deploySecret, cpu, ram, storage, keepAlive)

			_, err2 := clientset.CoreV1().Pods(namespaceName).Create(context.TODO(), podObject, metav1.CreateOptions{})
			if err2 != nil {
				r <- 1
				s.Stop()
				fmt.Println(err2.Error())
				s.Start()
			}

			for true {
				time.Sleep(time.Millisecond * 1000)
				pod, err := clientset.CoreV1().Pods(namespaceName).Get(context.TODO(), args[0], metav1.GetOptions{})
				if err != nil {
					r <- 1
					s.Stop()
					fmt.Println(err.Error())
					s.Start()
				}
				if pod.Status.Phase == "Running" {
					r <- 0
					break
				}
			}
		} else {
			r <- 1
			s.Stop()
			fmt.Println("Invalid target for tenant.")
			s.Start()
		}
	}()
	return r
}
