package get

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	v1 "k8s.io/api/core/v1"
	"kufast/clusterOperations"
	"kufast/tools"
)

// getCmd represents the get command
var getLogsCmd = &cobra.Command{
	Use:   "logs <podname>",
	Short: "Get the logs of a pod",
	Long:  `.`,
	Run: func(cmd *cobra.Command, args []string) {
		//Initial config block
		namespaceName, err := clusterOperations.GetTenantTargetNameFromCmd(cmd)
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

	getLogsCmd.Flags().StringP("target", "", "", "The name of the node to deploy the pod")
	getLogsCmd.Flags().StringP("tenant", "", "", "The name of the tenant to deploy the pod to")

}
