package tools

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func HandleError(err error, cmd *cobra.Command) {
	fmt.Println(err)
	_ = cmd.Help()
	os.Exit(1)
}
