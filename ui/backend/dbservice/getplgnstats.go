package dbservice

import (
	"strconv"

	hookDbService "github.com/aquasecurity/postee/dbservice"
	bolt "go.etcd.io/bbolt"
)

func GetPlgnStats() (r map[string]int, err error) {
	r = make(map[string]int)

	db, err := bolt.Open(hookDbService.DbPath, 0444, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(hookDbService.DbBucketOutputStats))
		if bucket == nil {
			return nil //no bucket - empty stats will be returned
		}

		c := bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			cnt, err := strconv.Atoi(string(v[:]))
			if err != nil {
				return err
			}

			r[string(k[:])] = cnt
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return r, nil
}
