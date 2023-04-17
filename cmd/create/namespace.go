package create

import (
	"github.com/spf13/cobra"
	"kufast/trackerFactory"
)

// createCmd represents the create command
var createNamespaceCmd = &cobra.Command{
	Use:   "namespace <name>..",
	Short: "Create one or more new namespaces for tenants",
	Long: `This command creates a new namespace for a tenant. You can select the name and set limits to the namespace.
Write multiple names to create multiple namespaces at once. This command will fail, if you do not have admin rights on the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {

		objectsCreated := len(args) * 2
		pw := trackerFactory.CreateProgressWriter(objectsCreated)
		for _, ns := range args {
			go trackerFactory.NewCreateNamespaceTracker(ns, cmd, pw)
		}
		trackerFactory.HandleTracking(pw, objectsCreated)
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
	createNamespaceCmd.Flags().StringP("limit-memory", "m", "4Gi", "Limit the RAM usage for this namespace")
	createNamespaceCmd.Flags().StringP("limit-cpu", "c", "2", "Limit the CPU usage for this namespace")
}
