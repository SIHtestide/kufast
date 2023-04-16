package list

import (
	"context"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kufast/tools"
	"os"
)

// listCmd represents the list command
var listUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "List all users in a namespace",
	Long: `List all users in your namespace. The overview contains the name of the user. And his permission within
the cluster. This command will fail, if you do not have admin rights on the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		//Initial config block
		ns, err := tools.GetNamespaceFromUserConfig(cmd)
		if err != nil {
			tools.HandleError(err, cmd)
		}

		clientset, _, err := tools.GetUserClient(cmd)
		if err != nil {
			tools.HandleError(err, cmd)
		}

		//execute request
		users, err := clientset.CoreV1().ServiceAccounts(ns).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			tools.HandleError(err, cmd)
		}

		//build table
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"NAME", "NAMESPACE", "CREATED AT"})
		for _, user := range users.Items {
			t.AppendRow(table.Row{user.Name, user.Namespace, user.CreationTimestamp})
			
		}
		t.AppendSeparator()
		t.Render()
	},
}

func init() {
	listCmd.AddCommand(listUsersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
