package dbservice

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"

	bolt "go.etcd.io/bbolt"
)

const (
	apiKeyName = "POSTEE_API_KEY"
)

func EnsureApiKey() error {
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

	newApiKey, err := generateApiKey(32)
	if err != nil {
		return err
	}

	err = dbInsert(db, DbBucketSharedConfig, []byte(apiKeyName), []byte(newApiKey))

	return err
}
func GetApiKey() (string, error) {
	var apiKey string = ""
	db, err := bolt.Open(DbPath, 0444, nil) //should be enough
	if err != nil {
		return "", err
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(DbBucketSharedConfig))
		if bucket == nil {
			return errors.New("no bucket") //no bucket
		}

		bytes := bucket.Get([]byte(apiKeyName))

		apiKey = string(bytes[:])
		return nil
	})
	if err != nil {
		return "", err
	}

	return apiKey, nil

}
func generateApiKey(length int) (string, error) {
	k := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return "", err
	}
	return hex.EncodeToString(k), nil
}
