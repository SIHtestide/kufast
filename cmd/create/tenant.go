package create

import (
	"context"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kufast/asyncOps"
	"kufast/objectFactory"
	"kufast/tools"
	"os"
	"time"
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

		//Configblock
		clientset, clientconfig, err := tools.GetUserClient(cmd)
		if err != nil {
			fmt.Println(err.Error())
		}

		//Activate spinner
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
		s.Prefix = "Creating Objects - Please wait!  "
		s.Start()

		_, err = clientset.CoreV1().ServiceAccounts("default").Create(context.TODO(), objectFactory.NewTenantUser(args[0], "default"), metav1.CreateOptions{})
		if err != nil {
			s.Stop()
			tools.HandleError(err, cmd)
		}

		//Read targets from Cobra
		targets, _ := cmd.Flags().GetStringArray("target")

		var createTargetOps []<-chan int32
		var targetResults []int32

		if targets != nil {
			for _, target := range targets {
				if tools.IsValidTarget(cmd, target, true) {
					createTargetOps = append(createTargetOps, asyncOps.CreateNamespace(args[0], target, cmd, s))
				} else {
					s.Stop()
					fmt.Println("Invalid target: " + target)
					s.Start()
				}

			}
		}

		//Ensure all operations are done
		for _, op := range createTargetOps {
			targetResults = append(targetResults, <-op)
		}

		user, err := clientset.CoreV1().ServiceAccounts("default").Get(context.TODO(), args[0]+"-user", metav1.GetOptions{})

		for i, res := range targetResults {
			if res == 0 {
				tools.AddTargetToTenant(cmd, targets[i], user)
			}
		}
		_, err = clientset.CoreV1().ServiceAccounts("default").Update(context.TODO(), user, metav1.UpdateOptions{})
		if err != nil {
			s.Stop()
			fmt.Println(err)
			s.Start()
		}
		tools.WriteNewUserYamlToFile(args[0], clientconfig, clientset, cmd, s)
		s.Stop()
		fmt.Println("Done!")
	},
}

func createTenantInteractive(cmd *cobra.Command, args []string) []string {

	return args
}

func init() {
	createCmd.AddCommand(createTenantCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
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
