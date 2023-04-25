package create

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"kufast/asyncOps"
	"kufast/tools"
	"os"
	"time"
)

// createCmd represents the create command
var createNamespaceCmd = &cobra.Command{
	Use:   "namespace <name>..",
	Short: "Create one or more new namespaces for tenants",
	Long: `This command creates a new namespace for a tenant. You can select the name and set limits to the namespace.
Write multiple names to create multiple namespaces at once. This command will fail, if you do not have admin rights on the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		isInteractive, _ := cmd.Flags().GetBool("interactive")
		if isInteractive {
			args = createNamespaceInteractive(cmd, args)
		}
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
		s.Prefix = "Creating Objects - Please wait!  "
		s.Start()
		var createOps []<-chan int32
		var results []int32

		for _, ns := range args {
			createOps = append(createOps, asyncOps.CreateNamespace(ns, cmd, s))
		}
		//Ensure all operations are done
		for _, op := range createOps {
			results = append(results, <-op)
		}
		s.Stop()
		fmt.Println("Done!")

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
	createNamespaceCmd.Flags().StringArrayP("users", "u", []string{}, "Usernames to create along with the namespace")
	createNamespaceCmd.Flags().StringP("output", "o", ".", "Folder to store the created client credentials. Mandatory, when defining -u")
	createNamespaceCmd.MarkFlagsRequiredTogether("output", "users")
	_ = createNamespaceCmd.MarkFlagDirname("output")
}

func createNamespaceInteractive(cmd *cobra.Command, args []string) []string {
	if len(args) == 0 {
		for true {
			namespaceName := tools.GetDialogAnswer("What should be the name for the new namespace?")
			args = append(args, namespaceName)
			continueNamespaces := tools.GetDialogAnswer("Do you want to add another namespace (yes/no)?")
			if continueNamespaces != "yes" {
				break
			}
		}
		limitCpu := tools.GetDialogAnswer("Which CPU limit do you want to set for the namespace(es)? (e.g. 1, 500m)")
		_ = cmd.Flags().Set("limit-cpu", limitCpu)
		limitMemory := tools.GetDialogAnswer("Which Memory limit do you want to set for the namespace(es)? (e.g. 100Mi, 5Gi)")
		_ = cmd.Flags().Set("limit-memory", limitMemory)

	}
	return args
}
