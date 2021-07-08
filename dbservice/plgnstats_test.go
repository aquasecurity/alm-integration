package dbservice

import (
	"os"
	"strconv"
	"testing"

	bolt "go.etcd.io/bbolt"
)

func TestRegisterPlgnInvctn(t *testing.T) {
	dbPathReal := DbPath
	defer func() {
		os.Remove(DbPath)
		DbPath = dbPathReal
	}()
	DbPath = "test_webhooks.db"
	expectedCnt := 3
	keyToTest := "test"
	for i := 0; i < expectedCnt; i++ {
		RegisterPlgnInvctn(keyToTest)
	}
	r, err := getPlgnStats()
	if err != nil {
		t.Fatal("error while getting value of API key")
	}
	if r[keyToTest] != expectedCnt {
		t.Errorf("Persisted count doesn't match expected. Expected %d, got %d\n", r[keyToTest], expectedCnt)
	}

}

func getPlgnStats() (r map[string]int, err error) {
	r = make(map[string]int)

	db, err := bolt.Open(DbPath, 0444, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(DbBucketOutputStats))
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
