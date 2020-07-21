package cmd

import (
	"os/exec"
)

// Exec is just wrapper for exec.Command to easier mocking
func Exec(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
