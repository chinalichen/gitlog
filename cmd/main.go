package main

import (
	"github.com/gin-gonic/gin"

	"github.com/chinalichen/gitlog/api"
)

func main() {
	r := gin.Default()

	gitLogApi := api.NewGitLogApiWrapper("./gitlog.db")
	gitLogApi.BindHandlers(r)

	// cmd := exec.Command("git", logprocess.GetGitLotArgs()...)
	// output, err := cmd.CombinedOutput()
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }
	// result, err := logprocess.Process(string(output))
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }
	// log.Print(result)
	r.Run(":8080")
}
