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

func GetDeploymentNamespace(cmd *cobra.Command) string {

	namespaceName, _ := GetNamespaceFromUserConfig(cmd)

	tenantName, _ := cmd.Flags().GetString("tenant")
	targetName, _ := cmd.Flags().GetString("target")

	if tenantName != "" && targetName != "" {
		namespaceName = tenantName + "-" + targetName
	} else if tenantName != "" {
		tenantName = GetTenantFromNamespace(namespaceName)
		namespaceName = tenantName + "-" + targetName
	} else if targetName != "" {
		namespaceName = tenantName + "-" + GetTenantDefaultTargetName(tenantName, cmd)
	}
	return namespaceName
}

func GetTenantFromNamespace(namespaceName string) string {
	return strings.Split(namespaceName, ("-"))[0]
}

func GetNamespaceFromUserConfig(cmd *cobra.Command) (string, error) {

	all, _ := cmd.Flags().GetBool("all-namespaces")
	if all {
		return "", nil
	}

	path, err := getKubeconfigPath(cmd)
	if err != nil {
		return "", err
	}

	ns, err := cmd.Flags().GetString("namespace")
	if err != nil {
		return "", err
	}
	if ns == "" {
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		loadingRules.Precedence[0] = path
		cfg, err := loadingRules.Load()
		if err != nil {
			return "", err
		} else if cfg.Contexts[cfg.CurrentContext] != nil {
			return cfg.Contexts[cfg.CurrentContext].Namespace, nil
		}

	} else {
		return ns, nil
	}

	return "", errors.New("Config not found or bad format.")
}

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
