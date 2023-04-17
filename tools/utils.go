package tools

import (
	"bufio"
	"context"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"os"
	"strings"
)

func HandleError(err error, cmd *cobra.Command) {
	fmt.Println(err)
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

func WriteNewUserYamlToFile(name string, namespaceName string, clientConfig *rest.Config, clientset *kubernetes.Clientset, out string, tracker progress.Tracker) {

	user, err := clientset.CoreV1().ServiceAccounts(namespaceName).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		fmt.Println(err)
	}

	secret, err := clientset.CoreV1().Secrets(namespaceName).Get(context.TODO(), user.Secrets[0].Name, metav1.GetOptions{})
	if err != nil {
		fmt.Println(err)
	}

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
			name: {
				Token: string(secret.Data["token"]),
			},
		},
		Contexts: map[string]*api.Context{
			"default-context": {
				Cluster:   "default-cluster",
				Namespace: namespaceName,
				AuthInfo:  name,
			},
		},
		CurrentContext: "default-context",
	}

	err = clientcmd.WriteToFile(newConfig, out+"/"+name+"-"+namespaceName+".kubeconfig")
	if err != nil {
		tracker.UpdateMessage("Unable to write config: " + err.Error())
		tracker.MarkAsErrored()
	}

}
