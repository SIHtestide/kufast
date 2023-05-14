// Package tools contains useful helper functions to reduce complexity and increase
// the maintainability of the code.
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
	"regexp"
	"strings"
	"syscall"
	"time"
)

// KUFAST_TENANT_DEFAULT_LABEL returns the name of the label for the default deployment namespace for a tenant
const KUFAST_TENANT_DEFAULT_LABEL = "kufast/default"

// KUFAST_TENANT_GROUPACCESS_LABEL returns the static part of the group access label of a tenant
const KUFAST_TENANT_GROUPACCESS_LABEL = "kufast.groupaccess/"

// KUFAST_TENANT_NODEACCESS_LABEL returns the static part of the node access label of a tenant
const KUFAST_TENANT_NODEACCESS_LABEL = "kufast.nodeaccess/"

// KUFAST_NODE_HOSTNAME_LABEL returns the Kubernetes default hostname label of a node
const KUFAST_NODE_HOSTNAME_LABEL = "kubernetes.io/hostname"

// KUFAST_NODE_GROUP_LABEL returns the static part of a group label that can be attached to a node
const KUFAST_NODE_GROUP_LABEL = "kufast.group/"

// KUFAST_TENANT_LABEL returns the default label for a tenant object
const KUFAST_TENANT_LABEL = "kufast/tenant"

// HandleError prints the error message given to it, prints the cobra commands help and exits the program
func HandleError(err error, cmd *cobra.Command) {
	fmt.Println("\n\n" + err.Error() + "\n\n")
	_ = cmd.Help()
	os.Exit(1)
}

// HandleErrorWithoutHelp prints the error message given to it and exits the program
func HandleErrorWithoutHelp(err error) {
	fmt.Println("\n\n" + err.Error() + "\n\n")
	os.Exit(1)
}

// GetDialogAnswer prints the given question to the user and expects and input to return
func GetDialogAnswer(question string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(question)
	fmt.Print(">> ")
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSuffix(answer, "\n")
	return answer
}

// GetPasswordAnswer prints the given question to the user and expects and input to return.
// The input is not shown on the command line
func GetPasswordAnswer(question string) string {
	fmt.Println(question)
	fmt.Print(">> ")
	password, _ := term.ReadPassword(int(syscall.Stdin))
	return strings.TrimSpace(string(password))
}

// WriteNewUserYamlToFile writes the credentials of a tenant to file. If the tenant has no tenant-target yet,
// the default namespace is set to the tenant-target user.
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
	} else {
		s.Stop()
		fmt.Println("Warning: No tenant-target specified! Consider to regenerate the tenants credentials after you created one" +
			" to avoid side effects!")
		s.Start()
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

func IsAlphaNumeric(s string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(s)
}
