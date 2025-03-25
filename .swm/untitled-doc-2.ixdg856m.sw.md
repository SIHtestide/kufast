---
title: Untitled doc (2)
---
# Introduction

This document will walk you through the process of creating a pod in our system. The purpose of this implementation is to facilitate the creation of Kubernetes pods using a command-line interface, allowing users to specify various parameters and options.

We will cover:

1. How the command for creating a pod is structured.
2. How the pod creation process is initiated and executed.
3. How the pod object is constructed and configured.
4. How the pod creation is handled asynchronously.

# Command structure

<SwmSnippet path="/cmd/create/pod.go" line="34">

---

The command for creating a pod is defined using the Cobra library. It specifies the command usage, description, and the function to execute when the command is run. This is where the user inputs the pod name and image.

```
// createPodCmd represents the create pod command
var createPodCmd = &cobra.Command{
	Use:   "pod <name> <image>",
	Short: "Create a new pod within a tenant-target",
	Long: `Creates a new pod within a tenant-target. A pod is like a shell for a container in Kuebrnetes. 
You need to specify the name and the image from which the pod should be created.
You can customize your deployment with the flags below or by using the interactive mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		isInteractive, _ := cmd.Flags().GetBool("interactive")
		if isInteractive {
			args = createPodInteractive(cmd)
		}
		if len(args) != 2 {
			tools.HandleError(errors.New(tools.ERROR_WRONG_NUMBER_ARGUMENTS), cmd)
		}
		s := tools.CreateStandardSpinner(tools.MESSAGE_CREATE_OBJECTS)
```

---

</SwmSnippet>

# Initiating pod creation

<SwmSnippet path="/cmd/create/pod.go" line="51">

---

The <SwmToken path="/cmd/create/pod.go" pos="41:1:1" line-data="	Run: func(cmd *cobra.Command, args []string) {">`Run`</SwmToken> function of the command checks for interactive mode and validates the arguments. It then starts a spinner to indicate the creation process and calls the <SwmToken path="/cmd/create/pod.go" pos="51:7:7" line-data="		res := clusterOperations.CreatePod(cmd, args)">`CreatePod`</SwmToken> function to handle the actual creation.

```
		res := clusterOperations.CreatePod(cmd, args)
		err := <-res
		s.Stop()
		if err != "" {
			tools.HandleError(errors.New(err), cmd)
		}
		fmt.Println(tools.MESSAGE_DONE)

	},
}
```

---

</SwmSnippet>

# Interactive pod creation

<SwmSnippet path="/cmd/create/pod.go" line="62">

---

If the user opts for interactive mode, the <SwmToken path="/cmd/create/pod.go" pos="62:2:2" line-data="// createPodInteractive is a helper function to create a pod interactively">`createPodInteractive`</SwmToken> function is invoked. This function prompts the user for various pod configuration options, such as CPU, memory, and storage limits, and sets these as command flags.

```
// createPodInteractive is a helper function to create a pod interactively
func createPodInteractive(cmd *cobra.Command) []string {
	fmt.Println(tools.MESSAGE_INTERACTIVE_IGNORE_INPUT)
	fmt.Println("This routine will create a new pod for you. Please note that your credentials must be available" +
		"according to our Readme and you need to hand over secrets, ports and an init command through the command line arguments" +
		"if you intend to use them. More information is available under 'kufast get pod --help'")
	//No arguments provided
	var args []string
	args = append(args, tools.GetDialogAnswer("Please enter the name, the pod should have in Kubernetes."))
	args = append(args, tools.GetDialogAnswer("Please enter the dockerimage, the pod should have."))
```

---

</SwmSnippet>

# Pod object construction

<SwmSnippet path="/objectFactory/deploymentObjects.go" line="26">

---

The <SwmToken path="/objectFactory/deploymentObjects.go" pos="32:2:2" line-data="// NewPod creates a new Kubernetes pod object based on several parameters.">`NewPod`</SwmToken> function in the <SwmToken path="/objectFactory/deploymentObjects.go" pos="24:2:2" line-data="package objectFactory">`objectFactory`</SwmToken> package is responsible for constructing the pod object. It takes several parameters, including the pod name, image, and resource limits, and returns a configured <SwmToken path="/objectFactory/deploymentObjects.go" pos="35:48:50" line-data="	attachedSecrets []string, deploySecret string, cpu string, ram string, storage string, shouldRestart bool, ports []int32, command []string) *v1.Pod {">`v1.Pod`</SwmToken> object.

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

# Configuring pod resources

<SwmSnippet path="/objectFactory/deploymentObjects.go" line="68">

---

Within the <SwmToken path="/objectFactory/deploymentObjects.go" pos="32:2:2" line-data="// NewPod creates a new Kubernetes pod object based on several parameters.">`NewPod`</SwmToken> function, resource limits for memory, CPU, and storage are set based on the provided parameters. This ensures that the pod is allocated the correct resources in the Kubernetes cluster.

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
```

---

</SwmSnippet>

# Handling pod creation asynchronously

<SwmSnippet path="/clusterOperations/pod.go" line="38">

---

The <SwmToken path="/clusterOperations/pod.go" pos="40:2:2" line-data="func CreatePod(cmd *cobra.Command, args []string) &lt;-chan string {">`CreatePod`</SwmToken> function in the <SwmToken path="/cmd/create/pod.go" pos="51:5:5" line-data="		res := clusterOperations.CreatePod(cmd, args)">`clusterOperations`</SwmToken> package handles the pod creation process asynchronously. It retrieves the necessary parameters from the command flags and uses the Kubernetes client to create the pod in the specified namespace.

```
// CreatePod creates a new pod as an async function. The input channel is closed, as soon as the operation
// completes. All parameters are drawn from the environment on the command line.
func CreatePod(cmd *cobra.Command, args []string) <-chan string {
	res := make(chan string)

	go func() {
		defer close(res)
		clientset, _, err := tools.GetUserClient(cmd)
		if err != nil {
			res <- err.Error()
			return
		}
```

---

</SwmSnippet>

# Monitoring pod creation

<SwmSnippet path="/clusterOperations/pod.go" line="81">

---

After initiating the pod creation, the function enters a loop to monitor the pod's status. It checks if the pod is running and returns an error message if the operation times out.

```
				time.Sleep(time.Millisecond * 1000)
				pod, err := clientset.CoreV1().Pods(namespaceName).Get(context.TODO(), args[0], metav1.GetOptions{})
				if err != nil {
					res <- err.Error()
					return
				}
				if pod.Status.Phase == "Running" {
					res <- ""
					break
				}
			}
```

---

</SwmSnippet>

This walkthrough provides an overview of the pod creation process, highlighting the key components and their roles in the system.

<SwmMeta version="3.0.0" repo-id="Z2l0aHViJTNBJTNBa3VmYXN0JTNBJTNBU0lIdGVzdGlkZQ==" repo-name="kufast"><sup>Powered by [Swimm](https://app.swimm.io/)</sup></SwmMeta>
