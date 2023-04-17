package tools

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func HandleError(err error, cmd *cobra.Command) {
	fmt.Println(err)
	_ = cmd.Help()
	os.Exit(1)
}

func GetDialogAnswer(question string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(question)
	fmt.Print(">>")
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSuffix(answer, "\n")
	return answer
}
