package executor

import (
	"fmt"
	"os/exec"
	"slices"

	"k8s.io/utils/ptr"

	"github.com/leemiyinghao/go-av1/internal/entity/task_flow"
)

// argsParse parses the arguments string into a slice of strings
// "" or ” are allowed in the arguments string
func argsParse(args string) []string {
	var parsedArgs []string
	var arg string
	var inQuote bool
	for _, r := range args {
		switch r {
		case ' ':
			if inQuote {
				arg += string(r)
			} else {
				parsedArgs = append(parsedArgs, arg)
				arg = ""
			}
		case '"':
			inQuote = !inQuote
		case '\'':
			inQuote = !inQuote
		default:
			arg += string(r)
		}
	}
	if arg != "" {
		parsedArgs = append(parsedArgs, arg)
	}
	return parsedArgs
}

func generateCommand(rawCommand string, file string) []string {
	var command []string
	for _, slice := range argsParse(rawCommand) {
		switch slice {
		case "$FILE":
			command = append(command, file)
		default:
			command = append(command, slice)
		}
	}
	if !slices.Contains(command, file) {
		command = append(command, file)
	}
	return command
}

func RunShellTask(task *task_flow.ShellTask, file string) (*string, *string, error) {
	// execute shell command and store the return code
	command := generateCommand(task.Command, file)
	executor := exec.Command(command[0], command[0:]...)
	var err error
	if err = executor.Start(); err != nil {
		return ptr.To("1"), &file, err
	}
	if err = executor.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return ptr.To(fmt.Sprintf("%d", exitError.ExitCode())), &file, err
		} else {
			return ptr.To("1"), &file, err
		}
	}
	return ptr.To("0"), &file, nil
}
