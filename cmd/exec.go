package cmd

import (
	"errors"
	"github.com/spf13/cobra"
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
	Short: "Exec into an existing container.",
	Long: `Exec into an existing container. With exec, you can access the command line of your
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
			tools.HandleError(errors.New("Too few or too many arguments provided."), cmd)
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

func init() {
	RootCmd.AddCommand(execCmd)

	// Here you will define your flags and configuration settings.
	execCmd.Flags().StringP("command", "c", "/bin/sh", "Set the command to be exec")
	execCmd.Flags().StringP("target", "", "", "The target for the secret (Needs to be the same as the pod using it.")
	execCmd.Flags().StringP("tenant", "", "", "The tenant owning the secret")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// execCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// execCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
