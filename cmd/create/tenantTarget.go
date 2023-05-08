// Package create contains the cmd functions for the creation of objects in kufast
package create

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"kufast/clusterOperations"
	"kufast/tools"
)

// createTenantTargetCmd represents the create tenant-target command
var createTenantTargetCmd = &cobra.Command{
	Use:   "tenant-target <target>..",
	Short: "Creates one or more new tenant-targets",
	Long: `Creates one or more new tenant-targets.
Tenant-targets will be attached to tenants and give them the ability to deploy pods to the target
until the specified resource limit is reached. Write multiple targets to create multiple tenant-targets at once. 
`,
	Run: func(cmd *cobra.Command, args []string) {

		isInteractive, _ := cmd.Flags().GetBool("interactive")
		if isInteractive {
			args = createTenantTargetInteractive(cmd)
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
			if !tools.IsAlphaNumeric(targetName) {
				s.Stop()
				fmt.Println(tools.CreateAlphaNumericError(tenantName))
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

		s.Stop()
		fmt.Println(tools.MESSAGE_DONE)

	},
}

// createTenantTargetInteractive is a helper function to create a tenant-target interactively
func createTenantTargetInteractive(cmd *cobra.Command) []string {
	fmt.Println(tools.MESSAGE_INTERACTIVE_IGNORE_INPUT)
	var args []string
	for true {
		namespaceName := tools.GetDialogAnswer("Please specify the target for the tenant-target.")
		args = append(args, namespaceName)
		continueNamespaces := tools.GetDialogAnswer("Do you want to add another target (y/N)?")
		if continueNamespaces != "y" {
			break
		}
	}

	limitCpu := tools.GetDialogAnswer("Which CPU limit do you want to set for the tenant-target(s)? (e.g. 1, 500m)")
	_ = cmd.Flags().Set("cpu", limitCpu)
	limitMemory := tools.GetDialogAnswer("Which memory limit do you want to set for the tenant-target(s)? (e.g. 100Mi, 5Gi)")
	_ = cmd.Flags().Set("memory", limitMemory)
	limitStorage := tools.GetDialogAnswer("Which storage limit do you want to set for the tenant-target(s)? (e.g. 100Mi, 5Gi)")
	_ = cmd.Flags().Set("storage", limitStorage)
	limitMinStorage := tools.GetDialogAnswer("Which minimum storage limit do you want to set for the tenant-target(s)? (e.g. 100Mi, 5Gi)")
	_ = cmd.Flags().Set("min-storage", limitMinStorage)
	limitPods := tools.GetDialogAnswer("Which pod limit do you want to set for the tenant-target(s)? (e.g. 100Mi, 5Gi)")
	_ = cmd.Flags().Set("pods", limitPods)

	return args
}

// init is a helper function from cobra to initialize the command. It sets all flags, standard values and documentation for this command.
func init() {
	createCmd.AddCommand(createTenantTargetCmd)

	//Set minimum values
	createTenantTargetCmd.Flags().StringP("memory", "", "1Gi", "Limit the RAM usage for the tenant-target(s)")
	createTenantTargetCmd.Flags().StringP("cpu", "", "500m", "Limit the CPU usage for the tenant-target(s)")
	createTenantTargetCmd.Flags().StringP("pods", "", "1", "Limit the Number of pods that can be created for the tenant-target(s)")
	createTenantTargetCmd.Flags().StringP("storage", "", "10Gi", "Limit the total storage for the tenant-target(s)")
	createTenantTargetCmd.Flags().StringP("storage-min", "", "1Gi", "Set the amount of storage, each pod must consume")

	//Tenant for the operation must be always specified
	createTenantTargetCmd.Flags().StringP("tenant", "t", "", "The tenant for the tenant-target(s).")
	_ = createTenantTargetCmd.MarkFlagRequired("tenant")

}
