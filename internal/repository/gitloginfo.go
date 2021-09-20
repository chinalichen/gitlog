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
	Name       string           `json:"name"`
	Status     GitLogInfoStatus `json:"status"`
	ModifiedAt time.Time        `json:"modifiedAt"`
}
