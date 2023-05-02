package tools

import (
	"bufio"
	"context"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"os"
	"strings"
	"syscall"
	"time"
)

const KUFAST_TENANT_DEFAULT_LABEL = "kufast/default"
const KUFAST_TENANT_GROUPACCESS_LABEL = "kufast.groupaccess/"
const KUFAST_TENANT_NODEACCESS_LABEL = "kufast.nodeaccess/"
const KUFAST_NODE_HOSTNAME_LABEL = "kubernetes.io/hostname"
const KUFAST_NODE_GROUP_LABEL = "kufast.group/"
const KUFAST_TENANT_LABEL = "kufast/tenant"

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

func WriteNewUserYamlToFile(tenantName string, cmd *cobra.Command, s *spinner.Spinner) error {

	clientset, clientConfig, err := GetUserClient(cmd)
	if err != nil {
		return err
	}

	tenant, errUser := clientset.CoreV1().ServiceAccounts("default").Get(context.TODO(), tenantName+"-user", metav1.GetOptions{})
	if errUser != nil {
		return err
	}
	secret, errSecret := clientset.CoreV1().Secrets("default").Get(context.TODO(), tenant.Secrets[0].Name, metav1.GetOptions{})
	if errSecret != nil {
		return err
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
	if tenant.ObjectMeta.Labels[KUFAST_TENANT_DEFAULT_LABEL] != "" {
		newConfig.Contexts["default-context"].Namespace = tenantName + "-" + tenant.ObjectMeta.Labels[KUFAST_TENANT_DEFAULT_LABEL]
	}

	err = clientcmd.WriteToFile(newConfig, out+"/"+tenantName+".kubeconfig")
	if err != nil {
		return err
	} else {
		s.Stop()
		fmt.Println("Config for tenant " + tenantName + " written to " + out + "/" + tenantName + ".kubeconfig")
		s.Start()
	}
	return nil
}

func CreateStandardSpinner(message string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	s.Prefix = message + "  "
	s.Start()

	return s
}
