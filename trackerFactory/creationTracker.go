package trackerFactory

import (
	"context"
	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kufast/objectFactory"
	"kufast/tools"
	"time"
)

func NewCreateNamespaceTracker(namespaceName string, cmd *cobra.Command, pw progress.Writer) {

	clientset, _, err := tools.GetUserClient(cmd)
	if err != nil {
		tools.HandleError(err, cmd)
	}

	ram, _ := cmd.Flags().GetString("limit-memory")
	cpu, _ := cmd.Flags().GetString("limit-cpu")

	tracker := progress.Tracker{Message: "Create Namespace base objects..", Total: 3, Units: progress.UnitsDefault}
	pw.AppendTracker(&tracker)

	_, err = clientset.CoreV1().Namespaces().Create(context.TODO(), objectFactory.NewNamespace(namespaceName), metav1.CreateOptions{})
	if err != nil {
		tracker.MarkAsErrored()
		tracker.Message = err.Error()
	}

	for true {
		newNamespace, err := clientset.CoreV1().Namespaces().Get(context.TODO(), namespaceName, metav1.GetOptions{})
		if err != nil {
			tracker.MarkAsErrored()
			tracker.Message = err.Error()
		}
		if newNamespace.Status.Phase == "Active" && !tracker.IsErrored() {
			tracker.SetValue(1)
			break
		}
		time.Sleep(time.Millisecond * 250)
	}

	if !tracker.IsErrored() {
		_, err = clientset.CoreV1().ResourceQuotas(namespaceName).Create(context.TODO(), objectFactory.NewResourceQuota(namespaceName, ram, cpu), metav1.CreateOptions{})
		if err != nil {
			tracker.MarkAsErrored()
			tracker.Message = err.Error()
		}
		tracker.SetValue(2)

		_, err = clientset.RbacV1().Roles(namespaceName).Create(context.TODO(), objectFactory.NewRole(namespaceName), metav1.CreateOptions{})
		if err != nil {
			tracker.MarkAsErrored()
			tracker.Message = err.Error()
		}

		tracker.SetValue(3)
		tracker.MarkAsDone()

		users, _ := cmd.Flags().GetStringArray("users")
		for _, user := range users {
			go NewCreateUserTracker(namespaceName, user, cmd, pw)
		}

	}
}

func NewCreateUserTracker(namespaceName string, userName string, cmd *cobra.Command, pw progress.Writer) {
	clientset, client, err := tools.GetUserClient(cmd)
	if err != nil {
		tools.HandleError(err, cmd)
	}

	tracker := progress.Tracker{Message: "Create User and Role Binding..", Total: 3, Units: progress.UnitsDefault}
	pw.AppendTracker(&tracker)

	_, err = clientset.CoreV1().ServiceAccounts(namespaceName).Create(context.TODO(), objectFactory.NewUser(userName, namespaceName), metav1.CreateOptions{})
	if err != nil {
		tracker.MarkAsErrored()
		tracker.Message = err.Error()
	}
	tracker.SetValue(1)
	_, err = clientset.RbacV1().RoleBindings(namespaceName).Create(context.TODO(), objectFactory.NewRoleBinding(userName, namespaceName), metav1.CreateOptions{})
	if err != nil {
		tracker.MarkAsErrored()
		tracker.Message = err.Error()
	}

	tracker.SetValue(2)
	out, _ := cmd.Flags().GetString("output")
	tools.WriteNewUserYamlToFile(userName, namespaceName, client, clientset, out, tracker)
	tracker.SetValue(3)
	tracker.MarkAsDone()

}

func NewCreatePodTracker(cmd *cobra.Command, pw progress.Writer, args []string) {
	clientset, _, err := tools.GetUserClient(cmd)
	if err != nil {
		tools.HandleError(err, cmd)
	}

	ram, _ := cmd.Flags().GetString("limit-memory")
	cpu, _ := cmd.Flags().GetString("limit-cpu")
	keepAlive, _ := cmd.Flags().GetBool("keep-alive")
	node, _ := cmd.Flags().GetString("target")
	secrets, _ := cmd.Flags().GetStringArray("secrets")

	namespaceName, _ := tools.GetNamespaceFromUserConfig(cmd)

	tracker := progress.Tracker{Message: "Create new Pod..", Total: 1, Units: progress.UnitsDefault}
	pw.AppendTracker(&tracker)

	podObject := objectFactory.NewPod(args[0], args[1], node, namespaceName, secrets, cpu, ram, keepAlive)

	_, err2 := clientset.CoreV1().Pods(namespaceName).Create(context.TODO(), podObject, metav1.CreateOptions{})
	if err2 != nil {
		tracker.MarkAsErrored()
		tracker.Message = err2.Error()
	}

	for true {
		time.Sleep(time.Millisecond * 1000)
		pod, err := clientset.CoreV1().Pods(namespaceName).Get(context.TODO(), args[0], metav1.GetOptions{})
		if err != nil {
			tracker.MarkAsErrored()
			tracker.Message = err.Error()
		}
		if pod.Status.Phase == "Running" && !tracker.IsErrored() {
			tracker.SetValue(1)
			break
		}
	}
}
