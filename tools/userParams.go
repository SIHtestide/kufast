package tools

import (
	"errors"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

func GetUserClient(path string) (*kubernetes.Clientset, *rest.Config, error) {
	config, err := clientcmd.BuildConfigFromFlags("", path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return clientset, config, err
	}

	return clientset, config, nil

}

func GetNamespaceFromUserConfig(path string) (string, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.Precedence[0] = path
	cfg, err := loadingRules.Load()
	if err != nil {
		return "", err
	} else if cfg.Contexts[cfg.CurrentContext] != nil {
		return cfg.Contexts[cfg.CurrentContext].Namespace, nil
	}

	return "", errors.New("Config not found or bad format.")
}
