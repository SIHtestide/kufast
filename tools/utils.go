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

const KUFAST_TENANT_DEFAULT_LABEL = "kufast/default"
const KUFAST_TENANT_GROUPACCESS_LABEL = "kufast.groupAccess/"
const KUFAST_TENANT_NODEACCESS_LABEL = "kufast.nodeAccess/"
const KUFAST_NODE_HOSTNAME_LABEL = "kubernetes.io/hostname/"
const KUFAST_NODE_GROUP_LABEL = "kufast.group/"
const KUFAST_TENANT_TARGET_ADMISSION_LABEL = "kufast.nodeAccess/"

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

func WriteNewUserYamlToFile(tenantName string, clientConfig *rest.Config, clientset *kubernetes.Clientset, cmd *cobra.Command, s *spinner.Spinner) {

	tenant, errUser := clientset.CoreV1().ServiceAccounts("default").Get(context.TODO(), tenantName+"-user", metav1.GetOptions{})
	if errUser != nil {
		s.Stop()
		fmt.Println(errUser)
		s.Start()
	}
	secret, errSecret := clientset.CoreV1().Secrets("default").Get(context.TODO(), tenant.Secrets[0].Name, metav1.GetOptions{})
	if errSecret != nil {
		s.Stop()
		fmt.Println(errSecret)
		s.Start()
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
			tenantName + "-user": {
				Token: string(secret.Data["token"]),
			},
		},
		Contexts: map[string]*api.Context{
			"default-context": {
				Cluster:   "default-cluster",
				Namespace: tenantName,
				AuthInfo:  tenantName + "-user",
			},
		},
		CurrentContext: "default-context",
	}
	if tenant.ObjectMeta.Labels["kufast/defaultTarget"] != "" {
		newConfig.Contexts["default-context"].Namespace = tenantName + "-" + tenant.ObjectMeta.Labels["kufast/defaultTarget"]
	}

	err := clientcmd.WriteToFile(newConfig, out+"/"+tenantName+".kubeconfig")
	if err != nil {
		s.Stop()
		fmt.Println("Unable to write config: " + err.Error())
		s.Start()
	} else {
		s.Stop()
		fmt.Println("Config for tenant " + tenantName + " written to " + out + "/" + tenantName + ".kubeconfig")
		s.Start()
	}
}
