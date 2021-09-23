package repository

import (
	"testing"
	"time"
)

func TestGetAndUpdateGitLogInfo(t *testing.T) {
	r := NewRepository("test.db")
	r.UpdateGitLogInfo(GitLogInfo{Name: "vscode", Status: GitLogInfoClone, ModifiedAt: time.Now()}, nil)
}
