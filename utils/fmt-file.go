package utils

import (
	"fmt"
	"os"
	"os/exec"
)

func FmtFile(path string) error {
	cmd := exec.Command("gofmt", "-w", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run gofmt: %v", err)
	}
	return nil
}
