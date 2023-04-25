package get

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	v1 "k8s.io/api/core/v1"
	"kufast/tools"
)

// getCmd represents the get command
var getLogsCmd = &cobra.Command{
	Use:   "logs <podname>",
	Short: "Get the logs of a pod",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//Initial config block
		namespaceName, err := tools.GetNamespaceFromUserConfig(cmd)
		if err != nil {
			tools.HandleError(err, cmd)
		}

		clientset, _, err := tools.GetUserClient(cmd)
		if err != nil {
			tools.HandleError(err, cmd)
		}
		count := int64(100)
		options := v1.PodLogOptions{
			Follow:    true,
			TailLines: &count,
		}

		//execute request
		log := clientset.CoreV1().Pods(namespaceName).GetLogs(args[0], &options)
		if err != nil {
			tools.HandleError(err, cmd)
		}

		podLogs, err := log.Stream(context.TODO())
		if err != nil {
			tools.HandleError(err, cmd)
		}

		defer podLogs.Close()

		for {
			buf := make([]byte, 2000)
			numBytes, err := podLogs.Read(buf)
			if numBytes == 0 {
				continue
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				tools.HandleError(err, cmd)
			}
			message := string(buf[:numBytes])
			fmt.Print(message)
		}

	},
}

func init() {
	getCmd.AddCommand(getLogsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
