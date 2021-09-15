package main

import (
	"fmt"
	"os/exec"
)

func main() {
	// tempDir := os.TempDir()
	cmd := exec.Command("git", "log --shortstat --pretty=format:'%h'")
	rc, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Print(err)
	}
	var out []byte
	_, err = rc.Read(out)
	if err != nil {
		fmt.Printf("command exit with %d", cmd.ProcessState.ExitCode())
	} else {
		fmt.Print(out)
	}
}
