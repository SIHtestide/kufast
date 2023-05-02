package create

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"kufast/clusterOperations"
	"kufast/tools"
)

// createCmd represents the create command
var createTenantTargetCmd = &cobra.Command{
	Use:   "tenant-target <target>..",
	Short: "Create one or more new namespaces for tenants",
	Long: `This command creates a new namespace for a tenant. You can select the name and set limits to the namespace.
Write multiple names to create multiple namespaces at once. This command will fail, if you do not have admin rights on the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {

		isInteractive, _ := cmd.Flags().GetBool("interactive")
		if isInteractive {
			args = createTenantTargetInteractive(cmd, args)
		}

		if len(args) < 1 {
			tools.HandleError(errors.New(tools.ERROR_WRONG_NUMBER_ARGUMENTS), cmd)
		}

		//Activate spinner
		s := tools.CreateStandardSpinner(tools.MESSAGE_CREATE_OBJECTS)

		tenantName, err := cmd.Flags().GetString("tenant")
		if err != nil {
			tools.HandleError(err, cmd)
		}

		var createTargetOps []<-chan string
		var targetResults []string

		for _, targetName := range args {
			err = clusterOperations.AddTargetToTenant(cmd, targetName, tenantName)
			if err != nil {
				s.Stop()
				fmt.Println(err)
				s.Start()
				continue
			}
			createTargetOps = append(createTargetOps, clusterOperations.CreateTenantTarget(tenantName, targetName, cmd))

		}

		//Ensure all operations are done
		for _, op := range createTargetOps {
			targetResults = append(targetResults, <-op)
		}

		for _, res := range targetResults {
			if res != "" {
				s.Stop()
				fmt.Println(res)
				s.Start()
			}
		}

		s.Stop()
		fmt.Println(tools.MESSAGE_DONE)

	},
}

func createTenantTargetInteractive(cmd *cobra.Command, args []string) []string {
	return args
}

func init() {
	createCmd.AddCommand(createTenantTargetCmd)

	createTenantTargetCmd.Flags().StringP("memory", "", "1Gi", "Limit the RAM usage for this namespace")
	createTenantTargetCmd.Flags().StringP("cpu", "", "500m", "Limit the CPU usage for this namespace")
	createTenantTargetCmd.Flags().StringP("pods", "", "1", "Limit the Number of pods that can be created in this namespace")
	createTenantTargetCmd.Flags().StringP("tenant", "t", "", "The tenant owning this namespace. Matching tenants will be connected.")
	createTenantTargetCmd.Flags().StringP("storage", "", "10Gi", "Limit the total storage in this namespace")
	createTenantTargetCmd.Flags().StringP("storage-min", "", "1Gi", "Set the amount of storage, each pod must consume")
	createTenantTargetCmd.Flags().StringArrayP("users", "u", []string{}, "Usernames to create along with the namespace")
	createTenantTargetCmd.Flags().StringP("output", "o", ".", "Folder to store the created client credentials. Mandatory, when defining -u")
	createTenantTargetCmd.MarkFlagsRequiredTogether("output", "users")
	_ = createTenantTargetCmd.MarkFlagDirname("output")
	_ = createTenantTargetCmd.MarkFlagRequired("tenant")

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
