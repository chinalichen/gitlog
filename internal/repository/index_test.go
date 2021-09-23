package repository

import (
	"testing"
	"time"
)

func TestGetAndUpdateGitLogInfo(t *testing.T) {
	r := NewRepository("test.db")
	info := GitLogInfo{Name: "vscode", Status: GitLogInfoClone, ModifiedAt: time.Now()}
	if err := r.UpdateGitLogInfo(info, nil); err != nil {
		t.Fatalf("update git log info error %v", err)
		return
	}
	
}
