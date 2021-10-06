package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/chinalichen/gitlog/api"
	"github.com/chinalichen/gitlog/internal/controller"
	"github.com/chinalichen/gitlog/internal/repository"
	"github.com/chinalichen/gitlog/pkg/gitprocess"
)

func main() {
	r := gin.Default()

	repo := repository.NewRepository("./gitlog.db")
	gp := gitprocess.NewGetProcessor(os.TempDir())
	manager := controller.NewLogManager(repo, gp)
	gitLogApi := api.NewGitLogApiWrapper(manager)

	r.GET("/git/info", gitLogApi.GetGitLogInfo)
	r.GET("/git/csv", gitLogApi.GetGitLogCSV)

	r.Run(":12389")
}
