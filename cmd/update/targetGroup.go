package update

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"kufast/clusterOperations"
	"kufast/tools"
)

// updateTargetGroupCmd represents the update target-group command
var updateTargetGroupCmd = &cobra.Command{
	Use:   "target-group <name> <nodes>",
	Short: "Create a new environment secret in this namespace",
	Long: `This command creates a new user and adds him to a namespace. You can select the namespace of the user.
Upon completion, the command yields the users credentials. This command will fail, if you do not have admin rights 
on the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 2 {
			tools.HandleError(errors.New(tools.ERROR_WRONG_NUMBER_ARGUMENTS), cmd)
		}

		s := tools.CreateStandardSpinner(tools.MESSAGE_UPDATE_OBJECTS)

		if clusterOperations.IsValidTarget(cmd, args[0], true) {
			err := clusterOperations.SetTargetGroupToNodes(args[0], args[:1], cmd)
			if err != nil {
				tools.HandleError(err, cmd)
			}
		}

		s.Stop()
		fmt.Println(tools.MESSAGE_DONE)

	},
}

// init is a helper function from cobra to initialize the command. It sets all flags, standard values and documentation for this command.
func init() {
	updateCmd.AddCommand(updateTargetGroupCmd)

}
