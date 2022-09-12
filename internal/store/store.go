package store

import (
	"go.etcd.io/bbolt"
)

var dbPath string

func InitDB(path string) error {
	dbPath = path

	db, err := bbolt.Open(dbPath, 0666, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	if err = db.Update(func(tx *bbolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte("users"))
		return err
	}); err != nil {
		return err
	}

	return nil
}
