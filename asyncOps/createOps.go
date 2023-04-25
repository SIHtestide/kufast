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

func CreateNamespace(namespaceName string, cmd *cobra.Command, s *spinner.Spinner) <-chan int32 {
	r := make(chan int32)

	go func() {
		defer close(r)

		clientset, _, err := tools.GetUserClient(cmd)
		if err != nil {
			tools.HandleError(err, cmd)
		}

		ram, _ := cmd.Flags().GetString("limit-memory")
		cpu, _ := cmd.Flags().GetString("limit-cpu")

		_, err = clientset.CoreV1().Namespaces().Create(context.TODO(), objectFactory.NewNamespace(namespaceName), metav1.CreateOptions{})
		if err != nil {
			r <- 1
			s.Stop()
			fmt.Println(err.Error())
			s.Start()
		}

		for true {
			newNamespace, err := clientset.CoreV1().Namespaces().Get(context.TODO(), namespaceName, metav1.GetOptions{})
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

		_, err = clientset.CoreV1().ResourceQuotas(namespaceName).Create(context.TODO(), objectFactory.NewResourceQuota(namespaceName, ram, cpu), metav1.CreateOptions{})
		if err != nil {
			r <- 1
			s.Stop()
			fmt.Println(err.Error())
			s.Start()
		}

		_, err = clientset.RbacV1().Roles(namespaceName).Create(context.TODO(), objectFactory.NewRole(namespaceName), metav1.CreateOptions{})
		if err != nil {
			r <- 1
			s.Stop()
			fmt.Println(err.Error())
			s.Start()
		}

		users, _ := cmd.Flags().GetStringArray("users")
		var createOps []<-chan int32
		var results []int32
		for _, user := range users {
			createOps = append(createOps, CreateUser(namespaceName, user, cmd, s))
		}
		//Ensure all operations are done
		for _, op := range createOps {
			results = append(results, <-op)
		}
		r <- 0
	}()
	return r

}

func CreateUser(namespaceName string, userName string, cmd *cobra.Command, s *spinner.Spinner) <-chan int32 {
	r := make(chan int32)

	go func() {
		defer close(r)
		clientset, client, err := tools.GetUserClient(cmd)
		if err != nil {
			r <- 1
			s.Stop()
			fmt.Println(err.Error())
			s.Start()
		}

		_, err = clientset.CoreV1().ServiceAccounts(namespaceName).Create(context.TODO(), objectFactory.NewUser(userName, namespaceName), metav1.CreateOptions{})
		if err != nil {
			r <- 1
			s.Stop()
			fmt.Println(err.Error())
			s.Start()
		}
		_, err = clientset.RbacV1().RoleBindings(namespaceName).Create(context.TODO(), objectFactory.NewRoleBinding(userName, namespaceName), metav1.CreateOptions{})
		if err != nil {
			r <- 1
			s.Stop()
			fmt.Println(err.Error())
			s.Start()
		}

		res := tools.WriteNewUserYamlToFile(userName, namespaceName, client, clientset, cmd, s)
		_ = <-res
		r <- 0

	}()
	return r
}

func NewCreatePodTracker(cmd *cobra.Command, s *spinner.Spinner, args []string) <-chan int32 {
	r := make(chan int32)

	go func() {
		clientset, _, err := tools.GetUserClient(cmd)
		if err != nil {
			r <- 1
			s.Stop()
			fmt.Println(err.Error())
			s.Start()
		}

		ram, _ := cmd.Flags().GetString("limit-memory")
		cpu, _ := cmd.Flags().GetString("limit-cpu")
		keepAlive, _ := cmd.Flags().GetBool("keep-alive")
		node, _ := cmd.Flags().GetString("target")
		secrets, _ := cmd.Flags().GetStringArray("secrets")

		namespaceName, _ := tools.GetNamespaceFromUserConfig(cmd)

		podObject := objectFactory.NewPod(args[0], args[1], node, namespaceName, secrets, cpu, ram, keepAlive)

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
	}()
	return r
}
