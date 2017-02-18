package tasks

import "os/exec"
import "strings"

// OSCommand - aka, exec.
func OSCommand(command string) (string, error) {
	var msg []byte
	var err error
	if strings.Contains(command, " ") {
		cmd := strings.SplitN(command, " ", 1)[0]
		args := strings.SplitAfterN(command, " ", 1)
		msg, err = exec.Command(cmd, args...).CombinedOutput()
	} else {
		msg, err = exec.Command("/bin/bash", "-c", command).CombinedOutput()
	}
	return string(msg), err
}
