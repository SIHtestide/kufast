// objectFactory builds new objects to be deployed to the cluster based on input parameters.
package objectFactory

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewPod creates a new Kubernetes pod object based on several parameters.
// Created objects only exist locally and need to be deployed to the cluster.
func NewPod(podName string, imageName string, namespaceName string,
	attachedSecrets []string, deploySecret string, cpu string, ram string, storage string, shouldRestart bool, ports []int32, command []string) *v1.Pod {

	var newPod *v1.Pod
	newPod = &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespaceName,
			Labels: map[string]string{
				"network": namespaceName,
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:    podName,
					Image:   imageName,
					Command: command,
					Resources: v1.ResourceRequirements{
						Limits:   v1.ResourceList{},
						Requests: v1.ResourceList{},
					},
					Ports: []v1.ContainerPort{},
					Env:   []v1.EnvVar{},
				},
			},
		},
		Status: v1.PodStatus{},
	}

	if ram != "" {
		qty, err := resource.ParseQuantity(ram)
		if err == nil {
			newPod.Spec.Containers[0].Resources.Limits["memory"] = qty
			newPod.Spec.Containers[0].Resources.Requests["memory"] = qty
		}
	}
	if cpu != "" {
		qty, err := resource.ParseQuantity(cpu)
		if err == nil {
			newPod.Spec.Containers[0].Resources.Limits["cpu"] = qty
			newPod.Spec.Containers[0].Resources.Requests["cpu"] = qty
		}
	}

	if storage != "" {
		qty, err := resource.ParseQuantity(storage)
		if err == nil {
			newPod.Spec.Containers[0].Resources.Limits["ephemeral-storage"] = qty
			newPod.Spec.Containers[0].Resources.Requests["ephemeral-storage"] = qty
		}
	}

	if shouldRestart {
		newPod.Spec.RestartPolicy = v1.RestartPolicyAlways
	}

	for _, port := range ports {
		containerPort := v1.ContainerPort{
			ContainerPort: port,
		}
		newPod.Spec.Containers[0].Ports = append(newPod.Spec.Containers[0].Ports, containerPort)
	}

	for _, secretName := range attachedSecrets {
		newPod.Spec.Containers[0].Env = append(newPod.Spec.Containers[0].Env, v1.EnvVar{
			Name: secretName,
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: secretName,
					},
					Key: "secret",
				},
			},
		})
	}

	if deploySecret != "" {
		newPod.Spec.ImagePullSecrets = []v1.LocalObjectReference{
			{
				Name: deploySecret,
			},
		}
	}

	return newPod

}

// NewSecret creates a new Kubernetes secret object based on several parameters.
// Created objects only exist locally and need to be deployed to the cluster.
func NewSecret(namespaceName string, secretName string, secretData string) *v1.Secret {
	return &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespaceName,
		},
		StringData: map[string]string{
			"secret": secretData,
		},
		Type: "Opaque",
	}
}

// NewDeploymentSecret creates a new Kubernetes secret object based on several parameters. This secret
// type can be used for Kubernetes deployments from private registries.
// Created objects only exist locally and need to be deployed to the cluster.
func NewDeploymentSecret(namespaceName string, secretName string, secretDataBase64 []byte) *v1.Secret {
	return &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespaceName,
		},
		Data: map[string][]byte{
			".dockerconfigjson": secretDataBase64,
		},
		Type: "kubernetes.io/dockerconfigjson",
	}
}
