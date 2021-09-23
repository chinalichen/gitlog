package repository

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
	"golang.org/x/xerrors"
)

const (
	GIT_LOG_INFO = "GitLogInfo"
	GIT_LOG_CSV  = "GitLogCSV"
)

type Repository struct {
	db *bolt.DB
}

func NewRepository(dbFile string) *Repository {
	r := &Repository{}

	if db, err := bolt.Open(dbFile, 0600, nil); err != nil {
		logrus.Fatal(err)
	} else {
		r.db = db
	}

	r.init()

	return r
}

func (r *Repository) init() {
	buckets := []string{GIT_LOG_INFO, GIT_LOG_CSV}

	// 生成 builtin 的bucket
	r.db.Update(func(tx *bolt.Tx) error {
		for _, bucketName := range buckets {
			_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
			if err != nil {
				return xerrors.Errorf("create %s bucket error %w", bucketName, err)
			}
		}
		return nil
	})
}

func (r *Repository) GetGitLogInfo(name string) (GitLogInfo, error) {
	var infoJson []byte
	info := GitLogInfo{}
	if err := r.db.View(func(tx *bolt.Tx) error {
		infoJson = tx.Bucket([]byte(GIT_LOG_INFO)).Get([]byte(name))
		return nil
	}); err != nil {
		return info, xerrors.Errorf("find info %s error: %w", name, err)
	}
	if len(infoJson) == 0 {
		return info, xerrors.Errorf("GetGitLogInfo %s lenth is 0", name)
	}
	json.Unmarshal(infoJson, &info)
	return info, nil
}

func (r *Repository) UpdateGitLogInfo(info GitLogInfo, batchTxIfNeeded *bolt.Tx) error {
	if batchTxIfNeeded != nil {
		return updateGitLogInfoImpl(info, batchTxIfNeeded)
	}
	return r.db.Update(func(tx *bolt.Tx) error {
		return updateGitLogInfoImpl(info, tx)
	})
}

func updateGitLogInfoImpl(info GitLogInfo, tx *bolt.Tx) error {
	infoJson, err := json.Marshal(info)
	if err != nil {
		return xerrors.Errorf("marshal GitLogInfo:%v error %w", info, err)
	}
	b := tx.Bucket([]byte(GIT_LOG_INFO))
	return b.Put([]byte(info.Name), infoJson)
}

func (r *Repository) GetGitLogCSV(name string) ([]byte, error) {
	var content []byte
	if err := r.db.View(func(tx *bolt.Tx) error {
		content = tx.Bucket([]byte(GIT_LOG_CSV)).Get([]byte(name))
		return nil
	}); err != nil {
		return nil, xerrors.Errorf("find csv %s error: %w", name, err)
	}
	if len(content) == 0 {
		return nil, xerrors.Errorf("GetGitLogCSV %s lenth is 0", name)
	}
	return content, nil
}

func (r *Repository) UpdateGitLogCSV(name string, content []byte, batchTxIfNeeded *bolt.Tx) error {
	if batchTxIfNeeded != nil {
		return updateGitLogCSVImpl(name, content, batchTxIfNeeded)
	}
	return r.db.Update(func(tx *bolt.Tx) error {
		return updateGitLogCSVImpl(name, content, tx)
	})
}

func updateGitLogCSVImpl(name string, content []byte, batchTx *bolt.Tx) error {
	b := batchTx.Bucket([]byte(GIT_LOG_CSV))
	return b.Put([]byte(name), content)
}

func (r *Repository) Batch(callback func(tx *bolt.Tx) error) {
	r.db.Batch(callback)
}

func (r *Repository) Close() {
	if r.db != nil {
		r.db.Close()
	}
}
