package controller

import "github.com/chinalichen/gitlog/internal/repository"

type LogManager struct {
	requests chan string
	r        *repository.Repository
}

func NewLogManager(dbFile string) *LogManager {
	lm := &LogManager{}
	lm.requests = make(chan string, 4)
	lm.r = repository.NewRepository(dbFile)
	return lm
}

func (lm *LogManager) GetGitLogInfo(name string) (repository.GitLogInfo, error) {
	return lm.r.GetGitLogInfo(name)
}
