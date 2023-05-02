package get

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"kufast/tools"
)

var getTenantCredsCmd = &cobra.Command{
	Use:   "tenant-creds <tenant>",
	Short: "Generate tenant credentials for specific tenant.",
	Long:  `Generate tenant credentials for specific user. Can only be used by admins.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			tools.HandleError(errors.New(tools.ERROR_WRONG_NUMBER_ARGUMENTS), cmd)
		}

		s := tools.CreateStandardSpinner(tools.MESSAGE_GET_OBJECTS)

		err := tools.WriteNewUserYamlToFile(args[0], cmd, s)
		if err != nil {
			s.Stop()
			tools.HandleError(err, cmd)
		}

		s.Stop()
		fmt.Println(tools.MESSAGE_DONE)

	},
}

func init() {
	getCmd.AddCommand(getTenantCredsCmd)
	getTenantCredsCmd.Flags().StringP("output", "o", ".", "Folder to store the created client credentials. Mandatory, when defining -u")
	_ = getTenantCredsCmd.MarkFlagRequired("output")

}
