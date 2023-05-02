package clusterOperations

import (
	"context"
	"errors"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kufast/objectFactory"
	"kufast/tools"
)

func CreateTenant(tenantName string, cmd *cobra.Command) error {

	//Configblock
	clientset, _, err := tools.GetUserClient(cmd)
	if err != nil {
		return err
	}

	_, err = clientset.CoreV1().ServiceAccounts("default").Create(context.TODO(), objectFactory.NewTenantUser(tenantName, "default"), metav1.CreateOptions{})
	if err != nil {
		return err
	}

	_, err = clientset.RbacV1().Roles("default").Create(context.TODO(), objectFactory.NewTenantDefaultRole(tenantName), metav1.CreateOptions{})
	if err != nil {
		return err
	}

	_, err = clientset.RbacV1().RoleBindings("default").Create(context.TODO(), objectFactory.NewTenantDefaultRoleBinding(tenantName), metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func DeleteTenant(tenantName string, cmd *cobra.Command) error {
	//Configblock
	clientset, _, err := tools.GetUserClient(cmd)
	if err != nil {
		return err
	}

	err = clientset.CoreV1().ServiceAccounts("default").Delete(context.TODO(), tenantName+"-user", metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	err = clientset.RbacV1().Roles("default").Delete(context.TODO(), tenantName+"-defaultrole", metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	err = clientset.RbacV1().RoleBindings("default").Delete(context.TODO(), tenantName+"-defaultrolebinding", metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}

func GetTenantNameFromCmd(cmd *cobra.Command) (string, error) {
	tenant, _ := cmd.Flags().GetString("tenant")
	if tenant == "" {
		namespaceName, err := tools.GetNamespaceFromUserConfig(cmd)
		if err != nil {
			return "", err
		}
		return tools.GetTenantFromNamespace(namespaceName), nil
	}
	return tenant, nil
}

func GetTenantFromCmd(cmd *cobra.Command) (*v1.ServiceAccount, error) {

	tenantName, err := GetTenantNameFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	//Configblock
	clientset, _, err := tools.GetUserClient(cmd)
	if err != nil {
		return nil, err
	}

	user, err := clientset.CoreV1().ServiceAccounts("default").Get(context.TODO(), tenantName+"-user", metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetTenantFromString(cmd *cobra.Command, tenantName string) (*v1.ServiceAccount, error) {

	//Configblock
	clientset, _, err := tools.GetUserClient(cmd)
	if err != nil {
		return nil, err
	}

	user, err := clientset.CoreV1().ServiceAccounts("default").Get(context.TODO(), tenantName+"-user", metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func UpdateTenantDefaultDeployTarget(newDefaultTarget string, cmd *cobra.Command) error {
	//Configblock
	clientset, _, err := tools.GetUserClient(cmd)
	if err != nil {
		return err
	}

	tenant, err := GetTenantFromCmd(cmd)
	if err != nil {
		return err
	}

	tenant.ObjectMeta.Labels[tools.KUFAST_TENANT_DEFAULT_LABEL] = newDefaultTarget
	_, err = clientset.CoreV1().ServiceAccounts("default").Update(context.TODO(), tenant, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil

}

func DeleteTargetFromTenant(targetName string, tenantName string, cmd *cobra.Command) error {
	if IsValidTenantTarget(cmd, targetName, tenantName, false) {
		clientset, _, err := tools.GetUserClient(cmd)
		if err != nil {
			return errors.New(err.Error())
		}

		target, err := GetTargetFromTargetName(cmd, targetName, tenantName, false)
		if err != nil {
			return errors.New(err.Error())
		}

		tenant, err := GetTenantFromString(cmd, tenantName)
		if err != nil {
			return errors.New(err.Error())
		}

		if target.AccessType == "node" {
			delete(tenant.ObjectMeta.Labels, tools.KUFAST_TENANT_NODEACCESS_LABEL+targetName)
		} else {
			delete(tenant.ObjectMeta.Labels, tools.KUFAST_TENANT_GROUPACCESS_LABEL+targetName)
		}
		_, err = clientset.CoreV1().ServiceAccounts("default").Update(context.TODO(), tenant, metav1.UpdateOptions{})
		if err != nil {
			return errors.New(err.Error())
		}

	} else {
		return errors.New("Not a valid target for this tenant: " + targetName)
	}

	return nil
}

func AddTargetToTenant(cmd *cobra.Command, targetName string, tenantName string) error {
	if IsValidTenantTarget(cmd, targetName, tenantName, true) {
		clientset, _, err := tools.GetUserClient(cmd)
		if err != nil {
			return errors.New(err.Error())
		}

		target, err := GetTargetFromTargetName(cmd, targetName, tenantName, true)
		if err != nil {
			return err
		}
		tenant, err := GetTenantFromString(cmd, tenantName)
		if err != nil {
			return err
		}
		if target.AccessType == "node" {
			tenant.ObjectMeta.Labels[tools.KUFAST_TENANT_NODEACCESS_LABEL+targetName] = "true"
		} else {
			tenant.ObjectMeta.Labels[tools.KUFAST_TENANT_GROUPACCESS_LABEL+targetName] = "true"
		}

		// Populate default label if possible
		if tenant.ObjectMeta.Labels[tools.KUFAST_TENANT_DEFAULT_LABEL] == "" {
			tenant.ObjectMeta.Labels[tools.KUFAST_TENANT_DEFAULT_LABEL] = targetName
		}
		_, err = clientset.CoreV1().ServiceAccounts("default").Update(context.TODO(), tenant, metav1.UpdateOptions{})
		if err != nil {
			return err
		}

		return nil
	}

	return errors.New("Invalid target!")
}

func GetTenantDefaultTargetNameFromCmd(cmd *cobra.Command) (string, error) {

	user, err := GetTenantFromCmd(cmd)
	if err != nil {
		return "", err
	}

	return user.ObjectMeta.Labels[tools.KUFAST_TENANT_DEFAULT_LABEL], nil
}
