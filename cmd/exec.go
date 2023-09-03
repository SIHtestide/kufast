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
package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"kufast/clusterOperations"
	"kufast/tools"
	"os"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec <pod>",
	Short: "Exec into an existing pod.",
	Long: `Exec into an existing pod. With exec, you can access the command line of your
container (provided one exists) and execute commands in the context of your container. This
will start an interactive CLI Session. To leave the container and get back to your normal
command line type "exit".`,
	Run: func(cmd *cobra.Command, args []string) {

		//Populate and set the command to be executed
		command, _ := cmd.Flags().GetString("command")

		//Populate namespace field
		namespaceName, err := clusterOperations.GetTenantTargetNameFromCmd(cmd)
		if err != nil {
			tools.HandleError(err, cmd)
		}

		clientset, config, err := tools.GetUserClient(cmd)
		if err != nil {
			tools.HandleError(err, cmd)
		}

		//Check that exactly one arg has been provided (the pod)
		if len(args) != 1 {
			tools.HandleError(errors.New(tools.ERROR_WRONG_NUMBER_ARGUMENTS), cmd)
		}

		//Set the command
		comm := []string{
			"sh",
			"-c",
			command,
		}

		req := clientset.CoreV1().RESTClient().Post().Resource("pods").Name(args[0]).
			Namespace(namespaceName).SubResource("exec")
		option := &v1.PodExecOptions{
			Command: comm,
			Stdin:   true,
			Stdout:  true,
			Stderr:  true,
			TTY:     true,
		}
		req.VersionedParams(
			option,
			scheme.ParameterCodec,
		)

		exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
		if err != nil {
			tools.HandleError(err, cmd)
		}
		err = exec.Stream(remotecommand.StreamOptions{
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		})
		if err != nil {
			tools.HandleError(err, cmd)
		}
	},
}

// init is a helper function from cobra to initialize the command. It sets all flags, standard values and documentation for this command.
func init() {
	RootCmd.AddCommand(execCmd)

	// Here you will define your flags and configuration settings.
	execCmd.Flags().StringP("command", "c", "/bin/sh", "Set the command for this operation")
	execCmd.Flags().StringP("target", "", "", tools.DOCU_FLAG_TARGET)
	execCmd.Flags().StringP("tenant", "", "", tools.DOCU_FLAG_TENANT)

}

func CreateExecDocs(linkH func(string) string) {
	out, err := os.Create("./kufast.wiki/exec.md")
	if err != nil {
		return
	}

	defer func() {
		err := out.Close()
		if err != nil {
			panic(err)
		}
	}()

	err = doc.GenMarkdownCustom(execCmd, out, linkH)
	if err != nil {
		panic(err)
	}

}
