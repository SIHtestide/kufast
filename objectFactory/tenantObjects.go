package objectFactory

import (
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	n1 "k8s.io/api/networking/v1"
	v12 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kufast/tools"
)

// NewNamespace creates a new Kubernetes namespace object based on several parameters.
// Created objects only exist locally and need to be deployed to the cluster.
func NewNamespace(tenantName string, target tools.Target, cmd *cobra.Command) *v1.Namespace {
	var newNamespace *v1.Namespace
	newNamespace = &v1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        tenantName + "-" + target.Name,
			Annotations: map[string]string{},
			Labels: map[string]string{
				tools.KUFAST_TENANT_LABEL: tenantName,
			},
		},
		Spec:   v1.NamespaceSpec{},
		Status: v1.NamespaceStatus{},
	}

	if target.AccessType == "node" {
		newNamespace.ObjectMeta.Annotations["scheduler.alpha.kubernetes.io/node-selector"] = tools.KUFAST_NODE_HOSTNAME_LABEL + "=" + target.Name
	} else {
		newNamespace.ObjectMeta.Annotations["scheduler.alpha.kubernetes.io/node-selector"] = tools.KUFAST_NODE_GROUP_LABEL + target.Name + "=true"
	}
	return newNamespace
}

// NewLimitRange creates a new Kubernetes LimitRange object based on several parameters.
// Created objects only exist locally and need to be deployed to the cluster.
func NewLimitRange(namespaceName string, minStorage string, storage string) *v1.LimitRange {
	var newRange *v1.LimitRange

	newRange = &v1.LimitRange{
		TypeMeta: metav1.TypeMeta{
			Kind:       "LimitRange",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      namespaceName + "-limitrange",
			Namespace: namespaceName,
		},
		Spec: v1.LimitRangeSpec{
			Limits: []v1.LimitRangeItem{
				{
					Type:           "Container",
					Min:            map[v1.ResourceName]resource.Quantity{},
					Max:            map[v1.ResourceName]resource.Quantity{},
					Default:        map[v1.ResourceName]resource.Quantity{},
					DefaultRequest: map[v1.ResourceName]resource.Quantity{},
				},
			},
		},
	}

	qty, err := resource.ParseQuantity(minStorage)
	if err == nil {
		newRange.Spec.Limits[0].Min["ephemeral-storage"] = qty
		newRange.Spec.Limits[0].Default["ephemeral-storage"] = resource.MustParse("1Gi")
		newRange.Spec.Limits[0].DefaultRequest["ephemeral-storage"] = resource.MustParse("1Gi")
	}
	qty, err = resource.ParseQuantity(storage)
	if err == nil {
		newRange.Spec.Limits[0].Max["ephemeral-storage"] = qty
	}

	return newRange
}

// NewResourceQuota creates a new Kubernetes ResourceQouta object based on several parameters.
// Created objects only exist locally and need to be deployed to the cluster.
func NewResourceQuota(namespace string, ram string, cpu string, storage string, pods string) *v1.ResourceQuota {
	var newQuota *v1.ResourceQuota
	newQuota = &v1.ResourceQuota{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ResourceQuota",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      namespace + "-limits",
			Namespace: namespace,
		},
		Spec: v1.ResourceQuotaSpec{
			Hard: v1.ResourceList{
				"secrets": resource.MustParse("100"),
			},
		},
	}
	//Set parameters only if available
	if ram != "" {
		qty, err := resource.ParseQuantity(ram)
		if err == nil {
			newQuota.Spec.Hard["limits.memory"] = qty
			newQuota.Spec.Hard["requests.memory"] = qty
		}
	}
	if cpu != "" {
		qty, err := resource.ParseQuantity(cpu)
		if err == nil {
			newQuota.Spec.Hard["limits.cpu"] = qty
			newQuota.Spec.Hard["requests.cpu"] = qty
		}
	}

	if pods != "" {
		qty, err := resource.ParseQuantity(pods)
		if err == nil {
			newQuota.Spec.Hard["pods"] = qty
		}
	}

	if storage != "" {
		qty, err := resource.ParseQuantity(storage)
		if err == nil {
			newQuota.Spec.Hard["requests.storage"] = qty
			newQuota.Spec.Hard["requests.ephemeral-storage"] = qty
			newQuota.Spec.Hard["limits.ephemeral-storage"] = qty
		}
	}
	return newQuota
}

// NewTenantUser creates a new Kubernetes ServiceAccount object based on several parameters.
// This is the basis user for a kufast tenant
// Created objects only exist locally and need to be deployed to the cluster.
func NewTenantUser(tenant string, namespaceName string) *v1.ServiceAccount {
	return &v1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      tenant + "-user",
			Namespace: namespaceName,
			Labels: map[string]string{
				tools.KUFAST_TENANT_LABEL:         tenant,
				tools.KUFAST_TENANT_DEFAULT_LABEL: "",
			},
		},
	}

}

// NewRole creates a new Kubernetes Role object based on several parameters.
// This role object is optimized for tenant targets.
// Created objects only exist locally and need to be deployed to the cluster.
func NewRole(namespaceName string) *v12.Role {
	return &v12.Role{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Role",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      namespaceName + "-role",
			Namespace: namespaceName,
		},
		Rules: []v12.PolicyRule{
			{
				APIGroups: []string{""},
				Verbs:     []string{"get", "list", "watch", "update", "delete", "create"},
				Resources: []string{"pods", "secrets", "pods/exec"},
			},
			{
				APIGroups: []string{""},
				Verbs:     []string{"get, list"},
				Resources: []string{"pods/log", "events"},
			},
		},
	}

}

// NewNetworkPolicy creates a new Kubernetes NetworkPolicy object based on several parameters.
// These network policies are preconfigured for Tenant Targets.
// Created objects only exist locally and need to be deployed to the cluster.
func NewNetworkPolicy(namespaceName string, tenant string) *n1.NetworkPolicy {

	return &n1.NetworkPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "NetworkPolicy",
			APIVersion: "networking.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      namespaceName + "-networkpolicy",
			Namespace: namespaceName,
		},
		Spec: n1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{},
			Ingress: []n1.NetworkPolicyIngressRule{
				{
					From: []n1.NetworkPolicyPeer{
						{
							NamespaceSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									tools.KUFAST_TENANT_LABEL: tenant,
								},
							},
						},
					},
				},
			},
			Egress: []n1.NetworkPolicyEgressRule{
				{
					To: []n1.NetworkPolicyPeer{
						{
							NamespaceSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									tools.KUFAST_TENANT_LABEL: tenant,
								},
							},
						},
					},
				},
			},
			PolicyTypes: []n1.PolicyType{
				"Ingress",
			},
		},
		Status: n1.NetworkPolicyStatus{},
	}

}

// NewTenantRolebinding creates a new Kubernetes RoleBinding object based on several parameters.
// This Role binding is preconfigured for the role binding of a tenant target role to a tenant.
// Created objects only exist locally and need to be deployed to the cluster.
func NewTenantRolebinding(namespaceName string, tenant string) *v12.RoleBinding {
	return &v12.RoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "RoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      namespaceName + "-" + tenant + "-binding",
			Namespace: namespaceName,
		},
		Subjects: []v12.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      tenant + "-user",
				Namespace: "default",
			},
		},
		RoleRef: v12.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     namespaceName + "-role",
		},
	}
}

// NewTenantDefaultRole creates a new Kubernetes RoleB object based on several parameters.
// This Role is preconfigured as kufast tenant standard role.
// Created objects only exist locally and need to be deployed to the cluster.
func NewTenantDefaultRole(tenantName string) *v12.Role {
	return &v12.Role{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Role",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      tenantName + "-defaultrole",
			Namespace: "default",
			Labels: map[string]string{
				tools.KUFAST_TENANT_LABEL: tenantName,
			},
		},
		Rules: []v12.PolicyRule{
			{
				Verbs:         []string{"get"},
				APIGroups:     []string{""},
				Resources:     []string{"serviceaccounts"},
				ResourceNames: []string{tenantName + "-user"},
			},
		},
	}

}

// NewTenantDefaultRoleBinding creates a new Kubernetes RoleB object based on several parameters.
// This Role binding is preconfigured for the role binding of the tenant default policy.
// Created objects only exist locally and need to be deployed to the cluster.
func NewTenantDefaultRoleBinding(tenantName string) *v12.RoleBinding {
	return &v12.RoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Role",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      tenantName + "-defaultrolebinding",
			Namespace: "default",
			Labels: map[string]string{
				tools.KUFAST_TENANT_LABEL: tenantName,
			},
		},
		Subjects: []v12.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      tenantName + "-user",
				Namespace: "default",
			},
		},
		RoleRef: v12.RoleRef{
			Kind: "Role",
			Name: tenantName + "-defaultrole",
		},
	}

}
