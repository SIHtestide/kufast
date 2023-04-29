package objectFactory

import (
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewNamespace(name string) *v1.Namespace {
	return &v1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec:   v1.NamespaceSpec{},
		Status: v1.NamespaceStatus{},
	}
}

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

func NewUser(name string, namespace string) *v1.ServiceAccount {
	return &v1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}

}

func NewRole(namespaceName string) *v12.Role {
	return &v12.Role{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Role",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      namespaceName + "-user",
			Namespace: namespaceName,
		},
		Rules: []v12.PolicyRule{
			{
				APIGroups: []string{""},
				Verbs:     []string{"get", "list", "watch", "update", "delete", "create"},
				Resources: []string{"pods", "secrets"},
			},
			{
				APIGroups: []string{""},
				Verbs:     []string{"get, list"},
				Resources: []string{"pods/log"},
			},
		},
	}

}

func NewRoleBinding(userName string, namespaceName string) *v12.RoleBinding {
	return &v12.RoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "RoleBinding",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      userName + "-" + namespaceName + "-role-binding",
			Namespace: namespaceName,
		},
		Subjects: []v12.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      userName,
				Namespace: namespaceName,
			},
		},
		RoleRef: v12.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     namespaceName + "-user",
		},
	}

}
