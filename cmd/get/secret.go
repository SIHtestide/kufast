package get

import (
	"errors"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"kufast/clusterOperations"
	"kufast/tools"
	"os"
)

// getCmd represents the get command
var getSecretCmd = &cobra.Command{
	Use:   "secret <secret>",
	Short: "Gain information about a deployed pod.",
	Long:  `Gain information about a deployed pod.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			tools.HandleError(errors.New(tools.ERROR_WRONG_NUMBER_ARGUMENTS), cmd)
		}

		s := tools.CreateStandardSpinner(tools.MESSAGE_GET_OBJECTS)

		secret, err := clusterOperations.GetSecret(args[0], cmd)
		if err != nil {
			tools.HandleError(err, cmd)
		}

		if secret.Data["secret"] == nil {
			err := errors.New("Error: This is not a valid secret")
			tools.HandleError(err, cmd)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"ATTRIBUTE", "VALUE"})
		t.AppendRow(table.Row{"Name", secret.Name})
		t.AppendRow(table.Row{"Namespace", secret.Namespace})
		t.AppendRow(table.Row{"Data", string(secret.Data["secret"])})

		s.Stop()

		t.AppendSeparator()
		t.Render()

	},
}

func init() {
	getCmd.AddCommand(getSecretCmd)

	getSecretCmd.Flags().StringP("target", "", "", "The name of the node to deploy the pod")
	getSecretCmd.Flags().StringP("tenant", "", "", "The name of the tenant to deploy the pod to")
}
