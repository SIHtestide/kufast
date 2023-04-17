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

	tracker := progress.Tracker{Message: "Create Namespace base objects..", Total: 1, Units: progress.UnitsDefault}
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

	}
}

func NewCreateUserTracker(namespaceName string, cmd *cobra.Command, pw progress.Writer) {
	
}
