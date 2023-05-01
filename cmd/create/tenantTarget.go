package create

import (
	"context"
	"errors"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kufast/asyncOps"
	"kufast/tools"
	"os"
	"time"
)

// createCmd represents the create command
var createTenantTargetCmd = &cobra.Command{
	Use:   "tenant-target <target>..",
	Short: "Create one or more new namespaces for tenants",
	Long: `This command creates a new namespace for a tenant. You can select the name and set limits to the namespace.
Write multiple names to create multiple namespaces at once. This command will fail, if you do not have admin rights on the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		//Configblock
		clientset, _, err := tools.GetUserClient(cmd)
		if err != nil {
			fmt.Println(err.Error())
		}

		isInteractive, _ := cmd.Flags().GetBool("interactive")
		if isInteractive {
			args = createTenantTargetInteractive(cmd, args)
		}

		if len(args) == 0 {
			tools.HandleError(errors.New("need to specify at least one target"), cmd)
		}

		user, err := clientset.CoreV1().ServiceAccounts("default").Get(context.TODO(), args[0]+"-user", metav1.GetOptions{})
		if err != nil {
			tools.HandleError(err, cmd)
		}

		//Activate spinner
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
		s.Prefix = "Creating Objects - Please wait!  "
		s.Start()

		tenant, _ := cmd.Flags().GetString("tenant")

		var createTargetOps []<-chan int32
		var targetResults []int32

		for _, targetName := range args {
			createTargetOps = append(createTargetOps, asyncOps.CreateNamespace(tenant, targetName, cmd, s))
		}
		//Ensure all operations are done
		for _, op := range createTargetOps {
			targetResults = append(targetResults, <-op)
		}

		//Add new Capabilities to user
		for i, res := range targetResults {
			if res == 0 {
				tools.AddTargetToTenant(cmd, args[i], user)
			}
		}

		_, err = clientset.CoreV1().ServiceAccounts("default").Update(context.TODO(), user, metav1.UpdateOptions{})
		if err != nil {
			s.Stop()
			fmt.Println(err)
			s.Start()
		}

		s.Stop()
		fmt.Println("Done!")

	},
}

func createTenantTargetInteractive(cmd *cobra.Command, args []string) []string {
	return args
}

func init() {
	createCmd.AddCommand(createTenantTargetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
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
