package cmd

import (
	"errors"
	"fmt"
	"os/exec"
)

// Exec is just wrapper for exec.Command to easier mocking
func Exec(name string, args ...string) (string, error) {
	c := exec.Command(name, args...)

	out, err := c.CombinedOutput()
	if err != nil {
		errMessage := fmt.Sprintf("%s: %s", err.Error(), out)
		return "", errors.New(errMessage)
	}

	return string(out), nil
}
