/*
MIT License

Copyright (c) 2023 Stefan Pawlowski

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
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

// getLogsCmd represents the get logs command
var getLogsCmd = &cobra.Command{
	Use:   "logs <podname>",
	Short: "Get the logs of a pod",
	Long:  `Get the logs of a pod.`,
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

// init is a helper function from cobra to initialize the command. It sets all flags, standard values and documentation for this command.
func init() {
	getCmd.AddCommand(getLogsCmd)

	getLogsCmd.Flags().StringP("target", "", "", tools.DOCU_FLAG_TARGET)
	getLogsCmd.Flags().StringP("tenant", "", "", tools.DOCU_FLAG_TENANT)

}
