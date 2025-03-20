---
title: User Command Interface
---
# Introduction

This document will walk you through the user command interface available in the kufast tool. The interface is designed to facilitate operations on Kubernetes clusters, allowing users to create, delete, update, list, and get information about various objects such as pods, secrets, tenants, and more.

We will cover:

&nbsp;

1. Get commands: How to retrieve information about existing objects.
2. List commands: How to list multiple objects within the cluster.
3. Delete commands: How to remove objects from the cluster.
4. Update commands: How to modify existing objects.
5. Create commands: How to add new objects to the cluster.
6. Exec commands: How to execute commands within a pod.

# Get commands

The get commands are used to retrieve detailed information about specific objects within the cluster. These commands provide insights into the attributes and status of the objects.

## Get deployment secret

<SwmSnippet path="/cmd/get/deployment-secret.go" line="35">

---

This command retrieves information about a <SwmToken path="/cmd/get/deployment-secret.go" pos="35:10:12" line-data="// getDeploySecretCmd represents the get deploy-secret command">`deploy-secret`</SwmToken>, including its name, <SwmToken path="/cmd/get/deployment-secret.go" pos="39:25:27" line-data="	Long:  `Gain information about a deploy-secret. Output includes name, tenant-target and the secret data.`,">`tenant-target`</SwmToken>, and secret data.

```
// getDeploySecretCmd represents the get deploy-secret command
var getDeploySecretCmd = &cobra.Command{
	Use:   "deploy-secret <secret>",
	Short: "Gain information about a deploy-secret.",
	Long:  `Gain information about a deploy-secret. Output includes name, tenant-target and the secret data.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			tools.HandleError(errors.New(tools.ERROR_WRONG_NUMBER_ARGUMENTS), cmd)
		}
```

---

</SwmSnippet>

## Get logs

<SwmSnippet path="/cmd/get/logs.go" line="36">

---

This command fetches the logs of a specified pod, allowing users to monitor its activity.

```
// getLogsCmd represents the get logs command
var getLogsCmd = &cobra.Command{
	Use:   "logs <podname>",
	Short: "Get the logs of a pod",
	Long:  `Get the logs of a pod.`,
	Run: func(cmd *cobra.Command, args []string) {
		//Initial config block
		namespaceName, err := clusterOperations.GetTenantTargetNameFromCmd(cmd)
		if err != nil {
			tools.HandleError(err, cmd)
		}
```

---

</SwmSnippet>

## Get pod

<SwmSnippet path="/cmd/get/pod.go" line="36">

---

This command provides detailed information about a deployed pod, including its status, resource limits, and events.

```
// getPodCmd represents the get pod command
var getPodCmd = &cobra.Command{
	Use:   "pod <pod>",
	Short: "Gain information about a deployed pod.",
	Long: `Gain information about a deployed pod. Output includes name, tenant-target, status, node, limits, image,
restart policy and IP-address`,
	Run: func(cmd *cobra.Command, args []string) {
```

---

</SwmSnippet>

## Get secret

<SwmSnippet path="/cmd/get/secret.go" line="35">

---

This command retrieves information about a secret, including its name, <SwmToken path="/cmd/get/secret.go" pos="39:23:25" line-data="	Long:  `Gain information about a secret. Output includes name, tenant-target and the secret data.`,">`tenant-target`</SwmToken>, and data.

```
// getSecretCmd represents the get secret command
var getSecretCmd = &cobra.Command{
	Use:   "secret <secret>",
	Short: "Gain information about a secret.",
	Long:  `Gain information about a secret. Output includes name, tenant-target and the secret data.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			tools.HandleError(errors.New(tools.ERROR_WRONG_NUMBER_ARGUMENTS), cmd)
		}
```

---

</SwmSnippet>

## Get tenant

<SwmSnippet path="/cmd/get/tenant.go" line="35">

---

This command provides information about a deployed tenant, including node and group access.

```
// getTenantCmd represents the get tenant command
var getTenantCmd = &cobra.Command{
	Use:   "tenant <tenant name>",
	Short: "Gain information about a deployed tenant.",
	Long:  `Gain information about a deployed tenant. Output includes name, node access and group access`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			tools.HandleError(errors.New(tools.ERROR_WRONG_NUMBER_ARGUMENTS), cmd)
		}
```

---

</SwmSnippet>

## Get tenant credentials

<SwmSnippet path="/cmd/get/tenantCredential.go" line="26">

---

This command generates credentials for a specific tenant, intended for admin use.

```
import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"kufast/tools"
)

// getTenantCredsCmd represents the get tenant-creds command
var getTenantCredsCmd = &cobra.Command{
	Use:   "tenant-creds <tenant>",
	Short: "Generate tenant credentials for specific tenant.",
	Long:  `Generate tenant credentials for specific user. Can only be used by admins.`,
	Run: func(cmd *cobra.Command, args []string) {
```

---

</SwmSnippet>

## Get tenant target

<SwmSnippet path="/cmd/get/tenantTarget.go" line="37">

---

This command provides information on a tenant target, including its resource limits and usage.

```
// getTenantTargetCmd represents the tenant-target command
var getTenantTargetCmd = &cobra.Command{
	Use:   "tenant-target <tenant-target>",
	Short: "Gain information on a tenant target.",
	Long:  `Gain information on a tenant target. Lists name, status, limits, their usage, and the number of included pods`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			tools.HandleError(errors.New(tools.ERROR_WRONG_NUMBER_ARGUMENTS), cmd)
		}
```

---

</SwmSnippet>

# List commands

List commands are used to display multiple objects within the cluster, providing an overview of their status and attributes.

## List pods

<SwmSnippet path="/cmd/list/pods.go" line="34">

---

This command lists all pods within a <SwmToken path="/cmd/list/pods.go" pos="37:15:17" line-data="	Short: &quot;List all pods in a tenant-target&quot;,">`tenant-target`</SwmToken>, showing their status and names.

```
// listPodsCmd represents the list pods command
var listPodsCmd = &cobra.Command{
	Use:   "pods",
	Short: "List all pods in a tenant-target",
	Long: `List all pods in yourtenant-target. The overview contains information about the status of your pod and its name.
To gain further information see the kubectl get pod command.`,
	Run: func(cmd *cobra.Command, args []string) {
```

---

</SwmSnippet>

## List secrets

<SwmSnippet path="/cmd/list/secrets.go" line="34">

---

This command lists all secrets within a <SwmToken path="/cmd/list/secrets.go" pos="37:15:17" line-data="	Short: &quot;List all secrets in a tenant-target&quot;,">`tenant-target`</SwmToken>, including their creation date.

```
// listSecretsCmd represents the list secrets command
var listSecretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "List all secrets in a tenant-target",
	Long: `List all secrets in a tenant-target. The overview contains information about the nameof the secret and its creation date
To gain further information see the kubectl get pod command.`,
	Run: func(cmd *cobra.Command, args []string) {
```

---

</SwmSnippet>

## List targets

<SwmSnippet path="/cmd/list/targets.go" line="34">

---

This command lists possible deployment targets for the current tenant, including nodes and groups.

```
// listTargetsCmd represents the list targets command
var listTargetsCmd = &cobra.Command{
	Use:   "targets",
	Short: "List possible targets for the current tenant.",
	Long:  `List possible targets for the current tenant. A target is a node or a group of nodes, a tenant can deploy to.`,
	Run: func(cmd *cobra.Command, args []string) {

		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			tools.HandleError(err, cmd)
		}
```

---

</SwmSnippet>

## List tenant targets

<SwmSnippet path="/cmd/list/tenantTargets.go" line="36">

---

This command lists all <SwmToken path="/cmd/list/tenantTargets.go" pos="36:10:12" line-data="// listTenantTargetsCmd represents the list tenant-targets command">`tenant-targets`</SwmToken> of a tenant, displaying their resource limits.

```
// listTenantTargetsCmd represents the list tenant-targets command
var listTenantTargetsCmd = &cobra.Command{
	Use:   "tenant-targets",
	Short: "List all tenant-targets of a tenant.",
	Long:  `List all tenant-targets of a tenant. The overview contains the limit information of each tenant target.`,
	Run: func(cmd *cobra.Command, args []string) {

		//Initial config block
		clientset, _, err := tools.GetUserClient(cmd)
		if err != nil {
			tools.HandleError(err, cmd)
		}
```

---

</SwmSnippet>

## List tenants

<SwmSnippet path="/cmd/list/tenants.go" line="36">

---

This command lists all tenants within the cluster, showing their namespaces and deployment targets.

```
// listTenantsCmd represents the list tenants command
var listTenantsCmd = &cobra.Command{
	Use:   "tenants",
	Short: "List all tenants in this cluster.",
	Long: `List all users in your namespace. The overview contains the name of the, the namespace where he is listed, the amount
of targets, this tenant can deploy to and the create date of this tenant.`,
	Run: func(cmd *cobra.Command, args []string) {
```

---

</SwmSnippet>

# Delete commands

Delete commands are used to remove objects from the cluster. These operations are irreversible, so caution is advised.

## Delete pod

<SwmSnippet path="/cmd/delete/pod.go" line="34">

---

This command deletes a specified pod, including its storage and logs.

```
// deletePodCmd represents the delete pod command
var deletePodCmd = &cobra.Command{
	Use:   "pod <pods>..",
	Short: "Delete the selected pod.",
	Long:  `Delete the selected pod including its storage. Please use with care! Deleted data cannot be restored.`,
	Run: func(cmd *cobra.Command, args []string) {

		//Check that exactly one arg has been provided (the namespace)
		if len(args) != 1 {
			tools.HandleError(errors.New(tools.ERROR_WRONG_NUMBER_ARGUMENTS), cmd)
		}
```

---

</SwmSnippet>

## Delete secret

<SwmSnippet path="/cmd/delete/secret.go" line="34">

---

This command deletes a secret from a <SwmToken path="/cmd/delete/secret.go" pos="37:15:17" line-data="	Short: &quot;Deletes a secret from a tenant-target.&quot;,">`tenant-target`</SwmToken>, applicable to both normal and <SwmToken path="/cmd/delete/secret.go" pos="39:34:36" line-data="Please use with care! Deleted data cannot be restored. Can be used for normal secrets and deploy-secrets.`,">`deploy-secrets`</SwmToken>.

```
// deleteSecretCmd represents the delete secret command
var deleteSecretCmd = &cobra.Command{
	Use:   "secret <secret>..",
	Short: "Deletes a secret from a tenant-target.",
	Long: `Deletes a secret from a tenant-target.
Please use with care! Deleted data cannot be restored. Can be used for normal secrets and deploy-secrets.`,
	Run: func(cmd *cobra.Command, args []string) {
		//Check that exactly one arg has been provided (the namespace)
		if len(args) < 1 {
			tools.HandleError(errors.New(tools.ERROR_WRONG_NUMBER_ARGUMENTS), cmd)
		}
```

---

</SwmSnippet>

## Delete target group

<SwmSnippet path="/cmd/delete/targetGroup.go" line="34">

---

This command deletes a <SwmToken path="/cmd/delete/targetGroup.go" pos="34:10:12" line-data="// deleteTargetGroupCmd represents the delete target-group command">`target-group`</SwmToken> from the cluster, affecting <SwmToken path="/cmd/list/tenantTargets.go" pos="36:10:12" line-data="// listTenantTargetsCmd represents the list tenant-targets command">`tenant-targets`</SwmToken> associated with it.

```
// deleteTargetGroupCmd represents the delete target-group command
var deleteTargetGroupCmd = &cobra.Command{
	Use:   "target-group <target-group>..",
	Short: "Deletes a target-group from the cluster.",
	Long: `Deletes a target-group from the cluster. This operation can only be executed by a cluster admin.
Please use with care! Tenant-targets pointing to these target-groups remain intact, but cannot deploy new pods.`,
	Run: func(cmd *cobra.Command, args []string) {
		//Check that exactly one arg has been provided (the namespace)
		if len(args) < 1 {
			tools.HandleError(errors.New("Too few arguments provided."), cmd)
		}
```

---

</SwmSnippet>

## Delete tenant

<SwmSnippet path="/cmd/delete/tenant.go" line="34">

---

This command deletes tenants and their associated objects, including <SwmToken path="/cmd/delete/tenant.go" pos="37:12:14" line-data="	Short: &quot;Delete tenants, their tenant-targets, pods, secrets and their credentials.&quot;,">`tenant-targets`</SwmToken> and credentials.

```
// deleteTenantCmd represents the delete tenant command
var deleteTenantCmd = &cobra.Command{
	Use:   "tenant <tenant>..",
	Short: "Delete tenants, their tenant-targets, pods, secrets and their credentials.",
	Long: `Delete tenants, their tenant-targets and their credentials. This operation can only be executed by a cluster admin.
Please use with care! Deleted data cannot be restored.`,
	Run: func(cmd *cobra.Command, args []string) {
		//Check that at least one tenant has been provided
		if len(args) < 1 {
			tools.HandleError(errors.New(tools.ERROR_WRONG_NUMBER_ARGUMENTS), cmd)
		}
```

---

</SwmSnippet>

## Delete tenant target

<SwmSnippet path="/cmd/delete/tenantTarget.go" line="39">

---

This command deletes <SwmToken path="/cmd/delete/tenantTarget.go" pos="42:7:9" line-data="	Short: &quot;Delete tenant-targets of a tenant including all pods and secrets in it.&quot;,">`tenant-targets`</SwmToken> of a tenant, including all pods and secrets within them.

```
// deleteTenantTargetCmd represents the delete tenant-target command
var deleteTenantTargetCmd = &cobra.Command{
	Use:   "tenant-target <tenant-target>",
	Short: "Delete tenant-targets of a tenant including all pods and secrets in it.",
	Long: `Delete tenant-targets of a tenant including all pods and secrets in it. This operation can only be executed by a cluster admin.
Please use with care! Deleted data cannot be restored.`,
	Run: func(cmd *cobra.Command, args []string) {
```

---

</SwmSnippet>

# Update commands

Update commands are used to modify existing objects within the cluster, adjusting their configurations and resource allocations.

## Update target group

<SwmSnippet path="/cmd/update/targetGroup.go" line="34">

---

This command updates the nodes within an existing target group, allowing reassignment of nodes.

```
// updateTargetGroupCmd represents the update target-group command
var updateTargetGroupCmd = &cobra.Command{
	Use:   "target-group <name> <nodes>..",
	Short: "Update the nodes on an existing target group.",
	Long: `Update the nodes on an existing target group. Specify all nodes that should be in the group after the reassignment. 
 Already existing pods on nodes will not be affected of this change.`,
	Run: func(cmd *cobra.Command, args []string) {
```

---

</SwmSnippet>

## Update tenant default

<SwmSnippet path="/cmd/update/tenantDefault.go" line="34">

---

This command sets a new default <SwmToken path="/cmd/update/tenantDefault.go" pos="37:13:15" line-data="	Short: &quot;Set a new default tenant-target for a tenant.&quot;,">`tenant-target`</SwmToken> for a tenant, ensuring the target is valid.

```
// updateTenantDefaultCmd represents the update tenant-default command
var updateTenantDefaultCmd = &cobra.Command{
	Use:   "tenant-default <newDefault>",
	Short: "Set a new default tenant-target for a tenant.",
	Long:  `Set a new default tenant-target for a tenant. The target must be valid for this tenant.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			tools.HandleError(errors.New(tools.ERROR_WRONG_NUMBER_ARGUMENTS), cmd)
		}
```

---

</SwmSnippet>

## Update tenant target

<SwmSnippet path="/cmd/update/tenantTarget.go" line="41">

---

This command updates the resource capabilities of a tenant target, including memory, CPU, and storage.

```
// updateTenantTargetCmd represents the update tenant-target command
var updateTenantTargetCmd = &cobra.Command{
	Use:   "tenant-target <tenant-target>",
	Short: "Update memory, CPU and storage capabilities of a tenant target.",
	Long: "Update memory, CPU and storage capabilities of a tenant target. " +
		"Also updates the role scheme to the latest version of kufast.",
	Run: func(cmd *cobra.Command, args []string) {
```

---

</SwmSnippet>

# Create commands

Create commands are used to add new objects to the cluster, enabling the deployment and management of resources.

## Create deployment secret

<SwmSnippet path="/cmd/create/deployment-secret.go" line="34">

---

This command creates a <SwmToken path="/cmd/create/deployment-secret.go" pos="34:10:12" line-data="// createDeploySecretCmd represents the create deploy-secret command">`deploy-secret`</SwmToken> in the specified target, necessary for accessing private registries.

```
// createDeploySecretCmd represents the create deploy-secret command
var createDeploySecretCmd = &cobra.Command{
	Use:   "deploy-secret name",
	Short: "Creates a deploy-secret in the specified target.",
	Long: `This command created a deploy-secret on the specified target. 
Deploy-secrets are required to create pods from private registries. This kind of secrets will be created
from dockerconfig files. A prerequisite to this command is, to store this file on your computer 
(default location is ~/.docker/config.json).
`,
	Run: func(cmd *cobra.Command, args []string) {
		isInteractive, _ := cmd.Flags().GetBool("interactive")
		if isInteractive {
			args = createDeploySecretInteractive(cmd)
		}
```

---

</SwmSnippet>

## Create pod

<SwmSnippet path="/cmd/create/pod.go" line="34">

---

This command creates a new pod within a <SwmToken path="/cmd/create/pod.go" pos="37:17:19" line-data="	Short: &quot;Create a new pod within a tenant-target&quot;,">`tenant-target`</SwmToken>, allowing customization through flags or interactive mode.

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

## Create secret

<SwmSnippet path="/cmd/create/secret.go" line="34">

---

This command creates a new environment secret within a namespace, requiring admin rights.

```
// createSecretCmd represents the create secret command
var createSecretCmd = &cobra.Command{
	Use:   "secret name",
	Short: "Create a new environment secret in this namespace",
	Long: `This command creates a new user and adds him to a namespace. You can select the namespace of the user.
Upon completion, the command yields the users credentials. This command will fail, if you do not have admin rights 
on the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		isInteractive, _ := cmd.Flags().GetBool("interactive")
		if isInteractive {
			args = createSecretInteractive(cmd)
		}
```

---

</SwmSnippet>

## Create target group

<SwmSnippet path="/cmd/create/targetGroup.go" line="34">

---

This command creates a <SwmToken path="/cmd/create/targetGroup.go" pos="34:10:12" line-data="// createTargetGroupCmd represents the create target-group command">`target-group`</SwmToken> within the cluster, assigning it to specified nodes.

```
// createTargetGroupCmd represents the create target-group command
var createTargetGroupCmd = &cobra.Command{
	Use:   "target-group <name> <nodes>",
	Short: "Create a target-group within the cluster",
	Long: `This command creates a new target-group and assigns it to the specified nodes.
Target-groups can be used to define a tenant-target that can deploy to a group of nodes,
instead of a single node.`,
	Run: func(cmd *cobra.Command, args []string) {
		isInteractive, _ := cmd.Flags().GetBool("interactive")
		if isInteractive {
			args = createTargetGroupInteractive()
		}
```

---

</SwmSnippet>

## Create tenant

<SwmSnippet path="/cmd/create/tenant.go" line="34">

---

This command creates one or more new tenants, providing credentials for cluster access.

```
// createTenantCmd represents the create tenant command
var createTenantCmd = &cobra.Command{
	Use:   "tenant <name>..",
	Short: "Creates one or more new tenants",
	Long: `Creates one or more new tenants.
A tenant is a separated entity that can work on your Kubernetes cluster. Pass the credentials for the tenant to a
partner that is supposed to work with your cluster. 
`,
	Run: func(cmd *cobra.Command, args []string) {
		isInteractive, _ := cmd.Flags().GetBool("interactive")
		if isInteractive {
			args = createTenantInteractive()
		}
		if len(args) < 1 {
			tools.HandleError(errors.New(tools.ERROR_WRONG_NUMBER_ARGUMENTS), cmd)
		}
```

---

</SwmSnippet>

## Create tenant target

<SwmSnippet path="/cmd/create/tenantTarget.go" line="34">

---

This command creates <SwmToken path="/cmd/create/tenantTarget.go" pos="37:15:17" line-data="	Short: &quot;Creates one or more new tenant-targets&quot;,">`tenant-targets`</SwmToken> for a tenant, enabling pod deployment within resource limits.

```
// createTenantTargetCmd represents the create tenant-target command
var createTenantTargetCmd = &cobra.Command{
	Use:   "tenant-target <target>..",
	Short: "Creates one or more new tenant-targets",
	Long: `Creates one or more new tenant-targets.
Tenant-targets will be attached to tenants and give them the ability to deploy pods to the target
until the specified resource limit is reached. Write multiple targets to create multiple tenant-targets at once. 
`,
	Run: func(cmd *cobra.Command, args []string) {
```

---

</SwmSnippet>

# Exec commands

Exec commands allow users to execute commands within a pod, providing access to the container's command line for interactive sessions.

## Exec pod

<SwmSnippet path="/cmd/delete/pod.go" line="65">

---

This command allows users to execute commands within a specified pod, starting an interactive CLI session.

```
			for _, res := range targetResults {
				if res != "" {
					s.Stop()
					fmt.Println(res)
					s.Start()
				}
			}
```

---

</SwmSnippet>

<SwmMeta version="3.0.0" repo-id="Z2l0aHViJTNBJTNBa3VmYXN0JTNBJTNBU0lIdGVzdGlkZQ==" repo-name="kufast"><sup>Powered by [Swimm](https://app.swimm.io/)</sup></SwmMeta>
