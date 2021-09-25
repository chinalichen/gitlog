package repository

import (
	"testing"
	"time"
)

func TestGetAndUpdateGitLogInfo(t *testing.T) {
	r := NewRepository("test.db")
	defer r.Close()
	info := GitLogInfo{URL: "vscode", Status: GitLogInfoClone, ModifiedAt: time.Now()}
	if err := r.UpdateGitLogInfo(info, nil); err != nil {
		t.Fatalf("update git log info error %v", err)
		return
	}
	fetchInfo, err := r.GetGitLogInfo(info.URL)
	if err != nil {
		t.Fatalf("get git log info error %v", err)
		return
	}
	if fetchInfo.URL != info.URL {
		t.Fatalf("expect %s, but got %s", info.URL, fetchInfo.URL)
		return
	}
	if fetchInfo.Status != info.Status {
		t.Fatalf("expect %d, but got %d", info.Status, fetchInfo.Status)
		return
	}
	if !fetchInfo.ModifiedAt.Equal(info.ModifiedAt) {
		t.Fatalf("expect %s, but got %s", info.ModifiedAt, fetchInfo.ModifiedAt)
		return
	}
}

func TestGetAndUpdateGitCSV(t *testing.T) {
	r := NewRepository("test.db")
	defer r.Close()
	content := []byte("col1,col2,col3\nv1,v2,v3")
	if err := r.UpdateGitLogCSV("vscode", content, nil); err != nil {
		t.Fatalf("get git log info error %v", err)
		return
	}
	csv, err := r.GetGitLogCSV("vscode")
	if err != nil {
		t.Fatalf("get git log csv error %v", err)
		return
	}
	if string(csv) != string(content) {
		t.Fatalf("expect %s, but got %s", content, csv)
		return
	}
}
