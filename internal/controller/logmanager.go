package controller

import (
	"sync"
	"time"

	"github.com/chinalichen/gitlog/internal/gitprocess"
	"github.com/chinalichen/gitlog/internal/repository"
	"github.com/google/uuid"
	"go.etcd.io/bbolt"
)

const (
	maxConcurrentProcesser = 4
)

type LogManager struct {
	requests chan string
	cloning  *sync.Map
	r        *repository.Repository
	gp       *gitprocess.GitProcessor
}

func NewLogManager(r *repository.Repository, gp *gitprocess.GitProcessor) *LogManager {
	lm := &LogManager{}
	lm.requests = make(chan string, maxConcurrentProcesser)
	lm.cloning = &sync.Map{}
	lm.r = r
	lm.gp = gp
	return lm
}

// func (lm *LogManager) Run() {
// 	for {
// 		select {
// 		case name := <-lm.requests:

// 		}
// 	}
// }

func (lm *LogManager) GetGitLogCSV(url string) ([]byte, error) {
	return lm.r.GetGitLogCSV(url)
}

func (lm *LogManager) GetGitLogInfo(url string) (repository.GitLogInfo, error) {
	if info, err := lm.r.GetGitLogInfo(url); err == nil {
		return info, nil
	}
	if err := lm.createGitLogInfo(url); err != nil {
		return repository.GitLogInfo{}, err
	}
	return lm.r.GetGitLogInfo(url)
}

func (lm *LogManager) createGitLogInfo(url string) error {
	isCloning, loaded := lm.cloning.LoadOrStore(url, true)
	if loaded && isCloning.(bool) {
		return nil
	}
	defer lm.cloning.Delete(url)
	info := repository.GitLogInfo{
		URL:        url,
		Path:       uuid.New().String(),
		Status:     repository.GitLogInfoCreate,
		ModifiedAt: time.Now(),
	}
	err := lm.r.UpdateGitLogInfo(info, nil)
	if err != nil {
		return err
	}

	go lm.pull(info.URL)

	return nil
}

func (lm *LogManager) pull(url string) error {
	info, err := lm.r.GetGitLogInfo(url)
	if err != nil {
		return err
	}
	info.Status = repository.GitLogInfoClone
	info.ModifiedAt = time.Now()
	if err := lm.r.UpdateGitLogInfo(info, nil); err != nil {
		return err
	}
	if err := lm.gp.Clone(info.Path, url); err != nil {
		return err
	}
	info2 := repository.GitLogInfo{
		URL:        url,
		Path:       info.Path,
		Status:     repository.GitLogInfoOutOfDate,
		ModifiedAt: time.Now(),
	}
	if err := lm.r.UpdateGitLogInfo(info2, nil); err != nil {
		return err
	}

	go lm.gitLog(url)

	return nil
}

func (lm *LogManager) gitLog(url string) error {
	info, err := lm.r.GetGitLogInfo(url)
	if err != nil {
		return err
	}
	csv, err := lm.gp.GitLog(info.Path)
	if err != nil {
		return err
	}
	info.Status = repository.GitLogInfoLatest
	lm.r.Batch(func(tx *bbolt.Tx) error {
		if err := lm.r.UpdateGitLogInfo(info, tx); err != nil {
			return err
		}
		if err := lm.r.UpdateGitLogCSV(info.URL, []byte(csv), tx); err != nil {
			return err
		}
		return nil
	})
	return nil
}
