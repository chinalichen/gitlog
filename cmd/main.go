package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	logprocess "github.com/chinalichen/gitlog/internal"
)

func main() {
	formatArgs := strings.Join([]string{"%h", "%p", "%an", "%ae", "%al", "%as", "%cN", "%ce", "%cl", "%cs", "%s"}, "#$@&")
	formatStr := fmt.Sprintf("--pretty=format:'%s'", formatArgs)

	log.Printf("will execute `%s`", strings.Join([]string{"git", "log", "--shortstat", formatStr}, " "))
	cmd := exec.Command("git", "log", "--shortstat", formatStr)
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Fatal(err)
	} else {
		result := logprocess.Process(string(output))
		log.Print(result)
	}
}
