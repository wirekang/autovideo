package exe

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func Run(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	fmt.Println(cmd.String())
	data, err := cmd.Output()

	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		return "", fmt.Errorf("stderr: %s", exitError.Stderr)
	}

	if err != nil {
		return "", err
	}

	fmt.Println(string(data))
	return strings.TrimSpace(string(data)), nil
}
