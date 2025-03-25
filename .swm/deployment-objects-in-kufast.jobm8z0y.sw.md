---
title: Deployment objects in Kufast
---
# Introduction

This document will walk you through the implementation of deployment objects in Kufast, specifically focusing on Kubernetes objects. The purpose is to understand how different Kubernetes deployment objects are generated and why certain design decisions were made.

We will cover:

1. How the <SwmToken path="/objectFactory/deploymentObjects.go" pos="32:2:2" line-data="// NewPod creates a new Kubernetes pod object based on several parameters.">`NewPod`</SwmToken> function constructs a Kubernetes Pod object.
2. How resource limits and requests are set for Pods.
3. How secrets and image pull secrets are integrated into Pods.
4. How the <SwmToken path="/objectFactory/deploymentObjects.go" pos="128:2:2" line-data="// NewSecret creates a new Kubernetes secret object based on several parameters.">`NewSecret`</SwmToken> and <SwmToken path="/objectFactory/deploymentObjects.go" pos="147:2:2" line-data="// NewDeploymentSecret creates a new Kubernetes secret object based on several parameters. This secret">`NewDeploymentSecret`</SwmToken> functions create Kubernetes Secret objects.

# Creating a Kubernetes Pod object

<SwmSnippet path="/objectFactory/deploymentObjects.go" line="26">

---

The <SwmToken path="/objectFactory/deploymentObjects.go" pos="32:2:2" line-data="// NewPod creates a new Kubernetes pod object based on several parameters.">`NewPod`</SwmToken> function is responsible for creating a new Kubernetes Pod object. This function takes several parameters, including the pod name, image name, namespace, and more. The created Pod object is configured locally and needs to be deployed to the Kubernetes cluster.

```
import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewPod creates a new Kubernetes pod object based on several parameters.
// Created objects only exist locally and need to be deployed to the cluster.
func NewPod(podName string, imageName string, namespaceName string,
	attachedSecrets []string, deploySecret string, cpu string, ram string, storage string, shouldRestart bool, ports []int32, command []string) *v1.Pod {
```

---

</SwmSnippet>

# Constructing the Pod object

<SwmSnippet path="/objectFactory/deploymentObjects.go" line="37">

---

The Pod object is constructed with essential metadata and specifications. This includes setting the Pod's name, namespace, and labels, as well as defining the container's image, command, and resource requirements.

```
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
```

---

</SwmSnippet>

# Setting resource limits and requests

<SwmSnippet path="/objectFactory/deploymentObjects.go" line="68">

---

Resource limits and requests for memory, CPU, and storage are set based on the provided parameters. This ensures that the Pod has the necessary resources allocated for its operation.

```
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
```

---

</SwmSnippet>

# Configuring restart policy

<SwmSnippet path="/objectFactory/deploymentObjects.go" line="91">

---

The restart policy for the Pod is configured based on the <SwmToken path="/objectFactory/deploymentObjects.go" pos="91:3:3" line-data="	if shouldRestart {">`shouldRestart`</SwmToken> parameter. If true, the Pod is set to always restart.

```
	if shouldRestart {
		newPod.Spec.RestartPolicy = v1.RestartPolicyAlways
	}
```

---

</SwmSnippet>

# Adding container ports

<SwmSnippet path="/objectFactory/deploymentObjects.go" line="95">

---

Container ports are added to the Pod based on the provided list of ports. This allows the Pod to expose the necessary ports for communication.

```
	for _, port := range ports {
		containerPort := v1.ContainerPort{
			ContainerPort: port,
		}
		newPod.Spec.Containers[0].Ports = append(newPod.Spec.Containers[0].Ports, containerPort)
	}
```

---

</SwmSnippet>

# Integrating secrets into the Pod

<SwmSnippet path="/objectFactory/deploymentObjects.go" line="102">

---

Secrets are integrated into the Pod's environment variables. This is done by appending each secret to the Pod's container environment, allowing the Pod to access sensitive information securely.

```
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
```

---

</SwmSnippet>

# Adding image pull secrets

<SwmSnippet path="/objectFactory/deploymentObjects.go" line="116">

---

Image pull secrets are added to the Pod if a <SwmToken path="/objectFactory/deploymentObjects.go" pos="116:3:3" line-data="	if deploySecret != &quot;&quot; {">`deploySecret`</SwmToken> is provided. This allows the Pod to pull images from private registries.

```
	if deploySecret != "" {
		newPod.Spec.ImagePullSecrets = []v1.LocalObjectReference{
			{
				Name: deploySecret,
			},
		}
	}
```

---

</SwmSnippet>

# Returning the constructed Pod

<SwmSnippet path="/objectFactory/deploymentObjects.go" line="124">

---

Finally, the constructed Pod object is returned, ready to be deployed to the Kubernetes cluster.

```
	return newPod

}
```

---

</SwmSnippet>

# Creating a Kubernetes Secret object

<SwmSnippet path="/objectFactory/deploymentObjects.go" line="128">

---

The <SwmToken path="/objectFactory/deploymentObjects.go" pos="128:2:2" line-data="// NewSecret creates a new Kubernetes secret object based on several parameters.">`NewSecret`</SwmToken> function creates a Kubernetes Secret object. This object is used to store sensitive data, such as passwords or tokens, in a secure manner.

```
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
```

---

</SwmSnippet>

# Creating a Kubernetes Deployment Secret object

<SwmSnippet path="/objectFactory/deploymentObjects.go" line="147">

---

The <SwmToken path="/objectFactory/deploymentObjects.go" pos="147:2:2" line-data="// NewDeploymentSecret creates a new Kubernetes secret object based on several parameters. This secret">`NewDeploymentSecret`</SwmToken> function creates a Kubernetes Secret object specifically for deployments from private registries. This secret type is used to store Docker configuration data in a base64-encoded format.

```
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
```

---

</SwmSnippet>

By understanding these functions and their implementations, you can effectively generate and manage Kubernetes deployment objects within the Kufast system.

<SwmMeta version="3.0.0" repo-id="Z2l0aHViJTNBJTNBa3VmYXN0JTNBJTNBU0lIdGVzdGlkZQ==" repo-name="kufast"><sup>Powered by [Swimm](https://app.swimm.io/)</sup></SwmMeta>
