package create

import (
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"kufast/tools"
	"kufast/trackerFactory"
	"os"
	"time"
)

// createCmd represents the create command
var createUserCmd = &cobra.Command{
	Use:   "user userName",
	Short: "Create a new user for a tenant",
	Long: `This command creates a new user and adds him to a namespace. You can select the namespace of the user.
Upon completion, the command yields the users credentials. This command will fail, if you do not have admin rights 
on the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		isInteractive, _ := cmd.Flags().GetBool("interactive")
		if isInteractive {
			args = createUserInteractive(args)
		}
		namespaceName, _ := tools.GetNamespaceFromUserConfig(cmd)

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
		s.Prefix = "Creating Objects - Please wait!  "
		s.Start()
		var createOps []<-chan int32
		var results []int32

		for _, user := range args {
			createOps = append(createOps, trackerFactory.NewCreateUserTracker(namespaceName, user, cmd, s))
		}
		//Ensure all operations are done
		for _, op := range createOps {
			results = append(results, <-op)
		}
		s.Stop()
	},
}

func init() {
	createCmd.AddCommand(createUserCmd)

	createUserCmd.Flags().StringP("output", "o", ".", "Folder to store the created client credentials. Mandatory, when defining -u")
	createUserCmd.MarkFlagRequired("output")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func createUserInteractive(args []string) []string {
	for true {
		args = append(args, tools.GetDialogAnswer("What name should the new user have?"))

		if tools.GetDialogAnswer("Do you want to create another user (yes/No)") != "yes" {
			break
		}
	}
	return args
}
