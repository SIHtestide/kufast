package trackerFactory

import (
	"context"
	"github.com/jedib0t/go-pretty/v6/progress"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
)

func NewCreateNamespaceTracker(namespace *v1.Namespace, quota *v1.ResourceQuota, clientset *kubernetes.Clientset, pw progress.Writer) {

	tracker := progress.Tracker{Message: "Create Namespace", Total: 1, Units: progress.UnitsDefault}
	pw.AppendTracker(&tracker)
	tracker2 := progress.Tracker{Message: "Create Limits", Total: 1, Units: progress.UnitsDefault}
	pw.AppendTracker(&tracker2)

	_, err := clientset.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
	if err != nil {
		tracker.MarkAsErrored()
		tracker.Message = err.Error()
	}

	for true {
		newNamespace, err := clientset.CoreV1().Namespaces().Get(context.TODO(), namespace.ObjectMeta.Name, metav1.GetOptions{})
		if err != nil {
			tracker.MarkAsErrored()
			tracker.Message = err.Error()
		}
		if newNamespace.Status.Phase == "Active" {
			tracker.MarkAsDone()
			break
		}
		time.Sleep(time.Millisecond * 250)
	}
	_, err = clientset.CoreV1().ResourceQuotas(namespace.ObjectMeta.Name).Create(context.TODO(), quota, metav1.CreateOptions{})
	if err != nil {
		tracker2.MarkAsErrored()
		tracker2.Message = err.Error()
	}

	tracker2.MarkAsDone()
}
