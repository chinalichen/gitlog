package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/chinalichen/gitlog/internal/controller"
	"github.com/chinalichen/gitlog/internal/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type getInfoQuery struct {
	URL string `form:"url"`
}

type GitLogApiWrapper struct {
	ctrl *controller.LogManager
}

func NewGitLogApiWrapper(lm *controller.LogManager) *GitLogApiWrapper {
	w := GitLogApiWrapper{}
	w.ctrl = lm
	return &w
}

func (w *GitLogApiWrapper) GetGitLogInfo(ctx *gin.Context) {
	q := getInfoQuery{}
	if err := ctx.Bind(&q); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	info, err := w.ctrl.GetGitLogInfo(q.URL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, info)
}

func (w *GitLogApiWrapper) GetGitLogCSV(ctx *gin.Context) {
	q := getInfoQuery{}
	if err := ctx.Bind(&q); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	info, err := w.ctrl.GetGitLogInfo(q.URL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if info.Status != repository.GitLogInfoLatest {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "refresh please"})
		return
	}
	csv, err := w.ctrl.GetGitLogCSV(q.URL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	content, err := utf8ToGb(csv)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	filePath := path.Join(os.TempDir(), fmt.Sprintf("%s.csv", info.Path))

	// if _, err := os.Stat(filePath); err != nil && os.IsExist(err) {
	// 	ctx.FileAttachment(filePath, getCsvName(info.URL))
	// 	return
	// } else if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
	// 	return
	// }

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if _, err := f.Write(content); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if err := f.Close(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctx.FileAttachment(filePath, getCsvName(info.URL))
}

func getCsvName(url string) string {
	nameList := strings.Split(url, "/")
	if len(nameList) > 0 {
		return fmt.Sprintf("%s.csv", nameList[len(nameList)-1])
	}
	return fmt.Sprintf("%s.csv", url)
}

func utf8ToGb(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GB18030.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
