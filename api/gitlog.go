package api

import (
	"net/http"

	"github.com/chinalichen/gitlog/internal/controller"
	"github.com/gin-gonic/gin"
)

type getInfoQuery struct {
	Name string `json:"name"`
}

type GitLogApiWrapper struct {
	ctrl *controller.LogManager
}

func NewGitLogApiWrapper(dbFile string) *GitLogApiWrapper {
	w := GitLogApiWrapper{}
	w.ctrl = controller.NewLogManager(dbFile)
	return &w
}

func (w *GitLogApiWrapper) BindHandlers(r *gin.Engine) {
	r.GET("/", w.getInfo)
}

func (w *GitLogApiWrapper) getInfo(ctx *gin.Context) {
	q := getInfoQuery{}
	if err := ctx.ShouldBindQuery(&q); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	info, err := w.ctrl.GetGitLogInfo(q.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, info)
}
