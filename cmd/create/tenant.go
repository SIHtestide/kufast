package create

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"kufast/clusterOperations"
	"kufast/tools"
)

// createCmd represents the create command
var createTenantCmd = &cobra.Command{
	Use:   "tenant <name>..",
	Short: "Create one or more new namespaces for tenants",
	Long: `This command creates a new namespace for a tenant. You can select the name and set limits to the namespace.
Write multiple names to create multiple namespaces at once. This command will fail, if you do not have admin rights on the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		isInteractive, _ := cmd.Flags().GetBool("interactive")
		if isInteractive {
			args = createTenantInteractive(cmd, args)
		}
		if len(args) < 1 {
			tools.HandleError(errors.New(tools.ERROR_WRONG_NUMBER_ARGUMENTS), cmd)
		}

		//Activate spinner
		s := tools.CreateStandardSpinner(tools.MESSAGE_CREATE_OBJECTS)

		for _, tenantName := range args {

			err := clusterOperations.CreateTenant(tenantName, cmd)
			if err != nil {
				tools.HandleError(err, cmd)
			}

			//Read targets from Cobra
			targets, _ := cmd.Flags().GetStringArray("target")

			var createTargetOps []<-chan string
			var targetResults []string

			if targets != nil {
				for _, targetName := range targets {
					err = clusterOperations.AddTargetToTenant(cmd, targetName, tenantName)
					if err != nil {
						s.Stop()
						fmt.Println(err)
						s.Start()
						err = clusterOperations.DeleteTargetFromTenant(targetName, tenantName, cmd)
						if err != nil {
							s.Stop()
							fmt.Println("Tenant is inconsistent. Please delete manually!")
						}
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
			}

			err = tools.WriteNewUserYamlToFile(tenantName, cmd, s)
			if err != nil {
				tools.HandleError(err, cmd)
			}

		}

		s.Stop()
		fmt.Println(tools.MESSAGE_DONE)
	},
}

func createTenantInteractive(cmd *cobra.Command, args []string) []string {

	return args
}

func init() {
	createCmd.AddCommand(createTenantCmd)

	createTenantCmd.Flags().StringP("memory", "", "1Gi", "Limit the RAM usage for this namespace")
	createTenantCmd.Flags().StringP("cpu", "", "500m", "Limit the CPU usage for this namespace")
	createTenantCmd.Flags().StringP("storage", "", "10Gi", "Limit the total storage in this namespace")
	createTenantCmd.Flags().StringP("storage-min", "", "1Gi", "Set the amount of storage, each pod must consume")
	createTenantCmd.Flags().StringP("pods", "", "1", "Limit the Number of pods that can be created in this namespace")

	createTenantCmd.Flags().StringArrayP("target", "", nil, "Deployment target for the namespace. Can be specified multiple times. Leave empty, for all targets")

	//Allow User definition
	createTenantCmd.Flags().StringP("output", "o", ".", "Folder to store the created client credentials. Mandatory, when defining -u")
	_ = createTenantCmd.MarkFlagDirname("output")
	_ = createTenantCmd.MarkFlagRequired("output")

	_ = createTenantCmd.MarkFlagRequired("tenant")

}
