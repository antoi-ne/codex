package store

import (
	"log"

	"go.etcd.io/bbolt"
)

var dbPath string

func InitDB(path string) {
	dbPath = path

	db, err := bbolt.Open(dbPath, 0666, nil)
	if err != nil {
		return
	}
	defer db.Close()

	if err = db.Update(func(tx *bbolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte("users"))
		return err
	}); err != nil {
		log.Fatalf("could not initiate database (%s)", err)
	}
}
