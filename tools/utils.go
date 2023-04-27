package tools

import (
	"bufio"
	"context"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"os"
	"strings"
	"syscall"
)

func HandleError(err error, cmd *cobra.Command) {
	fmt.Println("\n\n" + err.Error() + "\n\n")
	_ = cmd.Help()
	os.Exit(1)
}

func GetDialogAnswer(question string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(question)
	fmt.Print(">>")
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSuffix(answer, "\n")
	return answer
}

func GetPasswordAnswer(question string) string {
	fmt.Println(question)
	fmt.Print(">>")
	password, _ := term.ReadPassword(int(syscall.Stdin))
	return strings.TrimSpace(string(password))
}

func WriteNewUserYamlToFile(userName string, namespaceName string, clientConfig *rest.Config, clientset *kubernetes.Clientset, cmd *cobra.Command, s *spinner.Spinner) <-chan int32 {
	r := make(chan int32)

	go func() {
		defer close(r)
		user, errUser := clientset.CoreV1().ServiceAccounts(namespaceName).Get(context.TODO(), userName, metav1.GetOptions{})
		if errUser != nil {
			r <- 1
			s.Stop()
			fmt.Println(errUser.Error())
			s.Start()
			return
		}
		secret, errSecret := clientset.CoreV1().Secrets(namespaceName).Get(context.TODO(), user.Secrets[0].Name, metav1.GetOptions{})
		if errSecret != nil {
			r <- 1
			s.Stop()
			fmt.Println(errSecret)
			s.Start()
			return
		}

		out, _ := cmd.Flags().GetString("output")

		newConfig := api.Config{
			Kind:       "Config",
			APIVersion: "v1",
			Clusters: map[string]*api.Cluster{
				"default-cluster": {
					Server:                   clientConfig.Host,
					CertificateAuthorityData: secret.Data["ca.crt"],
				},
			},
			AuthInfos: map[string]*api.AuthInfo{
				userName: {
					Token: string(secret.Data["token"]),
				},
			},
			Contexts: map[string]*api.Context{
				"default-context": {
					Cluster:   "default-cluster",
					Namespace: namespaceName,
					AuthInfo:  userName,
				},
			},
			CurrentContext: "default-context",
		}

		err := clientcmd.WriteToFile(newConfig, out+"/"+userName+"-"+namespaceName+".kubeconfig")
		if err != nil {
			s.Stop()
			fmt.Println("Unable to write config: " + err.Error())
			s.Start()
			r <- 1
			return
		} else {
			s.Stop()
			fmt.Println("Config for user " + userName + " in namespace " + namespaceName + "written.")
			s.Start()
			r <- 0
			return
		}
	}()
	return r
}
