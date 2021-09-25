package repository

import "time"

type GitLogInfoStatus int8

const (
	GitLogInfoCreate GitLogInfoStatus = iota
	GitLogInfoClone
	GitLogInfoOutOfDate
	GitLogInfoLatest
)

type GitLogInfo struct {
	URL        string           `json:"url"`
	Path       string           `json:"path"`
	Status     GitLogInfoStatus `json:"status"`
	ModifiedAt time.Time        `json:"modifiedAt"`
}
