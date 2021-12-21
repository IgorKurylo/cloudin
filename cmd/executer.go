package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
)

func Executor(app string, stdout *string, stderror *string, args ...string) (error) {
	var stdOut bytes.Buffer
	var stdErr bytes.Buffer
	var cmd *exec.Cmd
	count := len(args)
	if count > 0 {
		cmd = exec.Command(app, args...)
	} else {
		cmd = exec.Command(app)
	}
	cmd.Stderr = &stdErr
	cmd.Stdout = &stdOut
	err := cmd.Run()
	fmt.Printf("Command: %q\n", cmd.Args)
	if err != nil {
		fmt.Println("Error: ", stdErr.String())
		*stderror = stdErr.String()
		return err
	}
	*stdout = stdOut.String()
	return nil

}

