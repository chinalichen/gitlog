package main

import (
	"log"
	"os/exec"

	logprocess "github.com/chinalichen/gitlog/internal"
)

func main() {
	cmd := exec.Command("git", logprocess.GetGitLotArgs()...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
		return
	}
	result, err := logprocess.Process(string(output))
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Print(result)
}
