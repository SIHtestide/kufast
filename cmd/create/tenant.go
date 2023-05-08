// Package create contains the cmd functions for the creation of objects in kufast
package create

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"kufast/clusterOperations"
	"kufast/tools"
)

// createTenantCmd represents the create tenant command
var createTenantCmd = &cobra.Command{
	Use:   "tenant <name>..",
	Short: "Creates one or more new tenants",
	Long: `Creates one or more new tenants.
A tenant is a separated entity that can work on your Kubernetes cluster. Pass the credentials for the tenant to a
partner that is supposed to work with your cluster. 
`,
	Run: func(cmd *cobra.Command, args []string) {
		isInteractive, _ := cmd.Flags().GetBool("interactive")
		if isInteractive {
			args = createTenantInteractive()
		}
		if len(args) < 1 {
			tools.HandleError(errors.New(tools.ERROR_WRONG_NUMBER_ARGUMENTS), cmd)
		}

		//Activate spinner
		s := tools.CreateStandardSpinner(tools.MESSAGE_CREATE_OBJECTS)

		for _, tenantName := range args {
			if !tools.IsAlphaNumeric(tenantName) {
				s.Stop()
				fmt.Println(tools.CreateAlphaNumericError(tenantName))
				s.Start()
				continue
			}

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
					if !tools.IsAlphaNumeric(targetName) {
						s.Stop()
						fmt.Println(tools.CreateAlphaNumericError(targetName))
						s.Start()
						continue
					}

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

// createTenantInteractive is a helper function to create a tenant interactively
func createTenantInteractive() []string {
	fmt.Println(tools.MESSAGE_INTERACTIVE_IGNORE_INPUT)
	var args []string
	args = append(args, tools.GetDialogAnswer("Please specify the name of the tenant."))
	for true {
		next := tools.GetDialogAnswer("Do you want to add another tenant? (y/N)")
		if next != "y" {
			break
		}
		args = append(args, tools.GetDialogAnswer("Please specify the name of the tenant."))
	}
	return args
}

// init is a helper function from cobra to initialize the command. It sets all flags, standard values and documentation for this command.
func init() {
	createCmd.AddCommand(createTenantCmd)

	createTenantCmd.Flags().StringP("memory", "", "1Gi", "Limit the RAM usage for this namespace")
	createTenantCmd.Flags().StringP("cpu", "", "500m", "Limit the CPU usage for this namespace")
	createTenantCmd.Flags().StringP("storage", "", "10Gi", "Limit the total storage in this namespace")
	createTenantCmd.Flags().StringP("storage-min", "", "1Gi", "Set the amount of storage, each pod must consume")
	createTenantCmd.Flags().StringP("pods", "", "1", "Limit the Number of pods that can be created in this namespace")

	createTenantCmd.Flags().StringArrayP("target", "", nil, "Deployment target for the namespace. Can be specified multiple times. Leave empty, for all targets")

	//Allow User definition
	createTenantCmd.Flags().StringP("output", "o", ".", "Folder to store the created client credentials.")
	_ = createTenantCmd.MarkFlagDirname("output")
	_ = createTenantCmd.MarkFlagRequired("output")

}
