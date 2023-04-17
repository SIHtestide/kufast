package objectFactory

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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

func NewResourceQuota(namespace string, ram string, cpu string) *v1.ResourceQuota {
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
			Hard: v1.ResourceList{},
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
	return newQuota
}

func NewUser(name string, namespace string) {

}

func NewRole(name string, namespace string) {

}

func NewRoleBinding(name string, namespace string) {

}

func NewUserYaml(name string, namespace string, clientset *kubernetes.Clientset) {

}
