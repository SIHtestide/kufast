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

	tracker := progress.Tracker{Message: "Create Namespace..", Total: 1, Units: progress.UnitsDefault}
	pw.AppendTracker(&tracker)
	tracker2 := progress.Tracker{Message: "Create Limits..", Total: 1, Units: progress.UnitsDefault}
	pw.AppendTracker(&tracker2)

	_, err = clientset.CoreV1().Namespaces().Create(context.TODO(), objectFactory.NewNamespace(namespaceName), metav1.CreateOptions{})
	if err != nil {
		tracker.MarkAsErrored()
		tracker2.MarkAsErrored()
		tracker.Message = err.Error()
		tracker2.Message = "Failed because of previous error"
	}

	for true {
		newNamespace, err := clientset.CoreV1().Namespaces().Get(context.TODO(), namespaceName, metav1.GetOptions{})
		if err != nil {
			tracker.MarkAsErrored()
			tracker2.MarkAsErrored()
			tracker.Message = err.Error()
			tracker2.Message = "Failed because of previous error"
		}
		if newNamespace.Status.Phase == "Active" && !tracker.IsErrored() {
			tracker.SetValue(1)
			tracker.MarkAsDone()
			break
		}
		time.Sleep(time.Millisecond * 250)
	}

	if !tracker.IsErrored() {
		_, err = clientset.CoreV1().ResourceQuotas(namespaceName).Create(context.TODO(), objectFactory.NewResourceQuota(namespaceName, ram, cpu), metav1.CreateOptions{})
		if err != nil {
			tracker2.MarkAsErrored()
			tracker2.Message = err.Error()
		}
		tracker2.SetValue(1)
		tracker2.MarkAsDone()
	}
}
