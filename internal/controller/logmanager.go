package controller

import (
	"sync"
	"time"

	"github.com/chinalichen/gitlog/internal/repository"
)

const (
	maxConcurrentProcesser = 4
)

type LogManager struct {
	requests chan string
	cloning  *sync.Map
	r        *repository.Repository
}

func NewLogManager(dbFile string) *LogManager {
	lm := &LogManager{}
	lm.requests = make(chan string, maxConcurrentProcesser)
	lm.cloning = &sync.Map{}
	lm.r = repository.NewRepository(dbFile)
	return lm
}

func (lm *LogManager) Run() {
	for {
		select {
		case name := <-lm.requests:

		}
	}
}

func (lm *LogManager) GetGitLogInfo(name string) (repository.GitLogInfo, error) {
	info, err := lm.r.GetGitLogInfo(name)
	if err == nil {
		return info, nil
	}
	lm.cloneGit(name)
}

func (lm *LogManager) cloneGit(name string) error {
	isCloning, loaded := lm.cloning.LoadOrStore(name, true)
	if loaded && isCloning.(bool) {
		return nil
	}
	defer lm.cloning.Delete(name)
	err := lm.r.UpdateGitLogInfo(repository.GitLogInfo{Name: name, Status: repository.GitLogInfoCreate, ModifiedAt: time.Now()}, nil)
	if err != nil {

	}
	lm.requests <- name
	return nil
}
