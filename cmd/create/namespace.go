package create

import (
	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kufast/tools"
	"kufast/trackerFactory"
	"time"
)

// createCmd represents the create command
var createNamespaceCmd = &cobra.Command{
	Use:   "namespace <name>",
	Short: "Create a new namespace for a tenant",
	Long: `This command creates a new namespace for a tenant. You can select the name and set limits to the namespace.
This command will fail, if you do not have admin rights on the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {

		clientset, _, err := tools.GetUserClient(cmd)
		if err != nil {
			tools.HandleError(err, cmd)
		}

		var newSpace *v1.Namespace
		newSpace = &v1.Namespace{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Namespace",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: args[0],
			},
			Spec:   v1.NamespaceSpec{},
			Status: v1.NamespaceStatus{},
		}

		var newQuota *v1.ResourceQuota
		newQuota = &v1.ResourceQuota{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ResourceQuota",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      args[0] + "-limits",
				Namespace: args[0],
			},
			Spec: v1.ResourceQuotaSpec{
				Hard: v1.ResourceList{
					"limits.cpu":      resource.MustParse("1"),
					"requests.cpu":    resource.MustParse("1"),
					"limits.memory":   resource.MustParse("1Gi"),
					"requests.memory": resource.MustParse("1Gi"),
				},
				Scopes:        nil,
				ScopeSelector: nil,
			},
			Status: v1.ResourceQuotaStatus{},
		}

		pw := progress.NewWriter()
		pw.SetNumTrackersExpected(2)
		pw.SetMessageWidth(100)
		pw.SetUpdateFrequency(time.Millisecond * 250)
		go pw.Render()
		go trackerFactory.NewCreateNamespaceTracker(newSpace, newQuota, clientset, pw)
		for !pw.IsRenderInProgress() {
			time.Sleep(time.Millisecond * 200)
		}
		for pw.IsRenderInProgress() {
			time.Sleep(time.Millisecond * 500)
			if pw.LengthDone() == 2 {
				pw.Stop()
			}
		}
	},
}

func init() {
	createCmd.AddCommand(createNamespaceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	createCmd.Flags().StringP("limit-ram", "r", "4G", "Limit the RAM usage for this namespace")
	createCmd.Flags().StringP("limit-cpu", "c", "2", "Limit the CPU usage for this namespace")
}
