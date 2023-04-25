package asyncOps

import (
	"context"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kufast/tools"
	"time"
)

func DeleteUser(userName string, cmd *cobra.Command, s *spinner.Spinner) <-chan int32 {
	r := make(chan int32)

	go func() {
		defer close(r)

		clientset, _, err := tools.GetUserClient(cmd)
		if err != nil {
			r <- 1
			s.Stop()
			fmt.Println(err.Error())
			s.Start()
		}

		namespaceName, _ := tools.GetNamespaceFromUserConfig(cmd)

		err = clientset.CoreV1().ServiceAccounts(namespaceName).Delete(context.TODO(), userName, v1.DeleteOptions{})
		if err != nil {
			r <- 1
			s.Stop()
			fmt.Println(err.Error())
			s.Start()
		}
		err = clientset.RbacV1().RoleBindings(namespaceName).Delete(context.TODO(), userName+"-"+namespaceName+"-role-binding", v1.DeleteOptions{})
		if err != nil {
			r <- 1
			s.Stop()
			fmt.Println(err.Error())
			s.Start()
		}

		for true {
			_, err := clientset.CoreV1().ServiceAccounts(namespaceName).Get(context.TODO(), userName, v1.GetOptions{})
			if err != nil {
				r <- 0
				break
			}
			time.Sleep(time.Millisecond * 250)
		}
		r <- 0

	}()
	return r
}
