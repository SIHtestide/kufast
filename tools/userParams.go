package tools

import (
	"errors"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"strings"
)

// GetUserClient creates an instance of clientset to communicate with the Kubernetes cluster
// based on the credentials the user entered when using this program.
func GetUserClient(cmd *cobra.Command) (*kubernetes.Clientset, *rest.Config, error) {
	var config *rest.Config
	var clientset *kubernetes.Clientset

	path, err := getKubeconfigPath(cmd)
	if err != nil {
		return clientset, config, err
	}

	config, err = clientcmd.BuildConfigFromFlags("", path)
	if err != nil {
		return clientset, config, err
	}

	// create the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return clientset, config, err
	}

	return clientset, config, nil
}

// GetTenantFromNamespace returns the tenants name from one of its namespaces by leveraging
// the namespace naming convention
func GetTenantFromNamespace(namespaceName string) string {
	return strings.Split(namespaceName, ("-"))[0]
}

// GetNamespaceFromUserConfig reads the userconfig of a user and returns the namespace
// specified in it.
func GetNamespaceFromUserConfig(cmd *cobra.Command) (string, error) {

	path, err := getKubeconfigPath(cmd)
	if err != nil {
		return "", err
	}

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.Precedence[0] = path
	cfg, err := loadingRules.Load()
	if err != nil {
		return "", err
	} else if cfg.Contexts[cfg.CurrentContext] != nil {
		return cfg.Contexts[cfg.CurrentContext].Namespace, nil
	} else {
		return "", errors.New("Config not found or bad format.")
	}

}

// getKubeconfigPath returns the path of the kubeconfig stored in a cobra command.
func getKubeconfigPath(cmd *cobra.Command) (string, error) {
	var kubeLoc string

	kubeLoc, err := cmd.Flags().GetString("kubeconfig")
	if err != nil {
		return "", err
	}
	if kubeLoc == "" {
		kubeLoc = homedir.HomeDir() + "/.kube/config"
	}

	return kubeLoc, nil
}
