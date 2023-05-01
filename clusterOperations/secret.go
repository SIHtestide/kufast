package clusterOperations

import (
	"context"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kufast/objectFactory"
	"kufast/tools"
	"os"
	"time"
)

func CreateDeploymentSecret(secretName string, cmd *cobra.Command) error {

	clientset, _, err := tools.GetUserClient(cmd)
	if err != nil {
		return err
	}

	fileName, err := cmd.Flags().GetString("input")
	if err != nil {
		return err
	}

	namespaceName, err := GetTenantTargetNameFromCmd(cmd)
	if err != nil {
		return err
	}

	creds, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	deploymentSecretObject := objectFactory.NewDeploymentSecret(namespaceName, secretName, creds)

	_, err = clientset.CoreV1().Secrets(namespaceName).Create(context.TODO(), deploymentSecretObject, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func CreateSecret(secretName string, secretData string, cmd *cobra.Command) error {
	//Default config block
	clientset, _, err := tools.GetUserClient(cmd)
	if err != nil {
		return err
	}

	//Get the namespace
	namespaceName, err := GetTenantTargetNameFromCmd(cmd)
	if err != nil {
		return err
	}

	//create secret object
	secretObject := objectFactory.NewSecret(namespaceName, secretName, secretData)

	//Push secret
	_, err = clientset.CoreV1().Secrets(namespaceName).Create(context.TODO(), secretObject, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func DeleteSecret(secretName string, cmd *cobra.Command) <-chan string {
	r := make(chan string)

	go func() {
		defer close(r)

		clientset, _, err := tools.GetUserClient(cmd)
		if err != nil {
			r <- err.Error()
			return
		}

		namespaceName, err := GetTenantTargetNameFromCmd(cmd)

		err = clientset.CoreV1().Secrets(namespaceName).Delete(context.TODO(), secretName, metav1.DeleteOptions{})
		if err != nil {
			r <- err.Error()
			return
		}

		for true {
			_, err := clientset.CoreV1().Secrets(namespaceName).Get(context.TODO(), secretName, metav1.GetOptions{})
			if err != nil {
				r <- err.Error()
				return
			}
			time.Sleep(time.Millisecond * 250)
		}
		r <- ""

	}()
	return r
}
