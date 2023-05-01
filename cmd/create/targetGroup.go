package create

import (
	"errors"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"kufast/tools"
	"os"
	"time"
)

// createCmd represents the create command
var createTargetGroupCmd = &cobra.Command{
	Use:   "target-group <name> <nodes>",
	Short: "Create a new environment secret in this namespace",
	Long: `This command creates a new user and adds him to a namespace. You can select the namespace of the user.
Upon completion, the command yields the users credentials. This command will fail, if you do not have admin rights 
on the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 2 {
			tools.HandleError(errors.New("too few arguments provided"), cmd)
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
		s.Prefix = "Creating Objects - Please wait!  "
		s.Start()

		err := tools.SetTargetGroupToNodes(args[0], args[:1], cmd)
		if err != nil {
			tools.HandleError(err, cmd)
		}

		s.Stop()
		fmt.Println("Complete!")

	},
}

func init() {
	createCmd.AddCommand(createTargetGroupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
