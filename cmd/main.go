package main

import (
	"log"
	"os/exec"

	"github.com/gin-gonic/gin"

	"github.com/chinalichen/gitlog/api"
	logprocess "github.com/chinalichen/gitlog/internal/logprocess"
)

func main() {
	r := gin.Default()

	api.GitLog(r)

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
