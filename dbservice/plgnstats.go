package dbservice

import (
	"strconv"

	bolt "go.etcd.io/bbolt"
)

func RegisterPlgnInvctn(name string) error {
	mutex.Lock()
	defer mutex.Unlock()

	db, err := bolt.Open(DbPath, 0666, nil)
	if err != nil {
		return err
	}
	defer db.Close()
	err = Init(db, DbBucketOutputStats)
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(DbBucketOutputStats))
		var i int
		v := bucket.Get([]byte(name))

		if v == nil {
			i = 0
		} else {
			i, err = strconv.Atoi(string(v[:]))
		}

		i++
		nwv := strconv.Itoa(i)

		err = bucket.Put([]byte(name), []byte(nwv))
		return err
	})

	return err
}
