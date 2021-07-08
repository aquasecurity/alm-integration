package dbservice

import (
	"errors"
	"os"
	"testing"

	"github.com/aquasecurity/postee/data"
	"go.etcd.io/bbolt"
)

var tests = []struct {
	caseDesc        string
	errPrvdr        func() error
	initIsNotCalled bool
}{
	{
		caseDesc: "EnsureApiKey",
		errPrvdr: func() error {
			return EnsureApiKey()
		},
	},
	{
		caseDesc: "GetApiKey",
		errPrvdr: func() error {
			_, err := GetApiKey()
			return err
		},
		initIsNotCalled: true,
	},
	{
		caseDesc: "RegisterPlgnInvctn",
		errPrvdr: func() error {
			return RegisterPlgnInvctn("some-key")
		},
	},
	{
		caseDesc: "HandleCurrentInfo",
		errPrvdr: func() error {
			_, _, err := HandleCurrentInfo(&data.ScanImageInfo{})
			return err
		},
	},
	{
		caseDesc: "AggregateScans",
		errPrvdr: func() error {
			_, err := AggregateScans("", map[string]string{}, 1, false)
			return err
		},
	},
}

func TestInvalidDbPath(t *testing.T) {

	dbPathReal := DbPath
	defer func() {
		os.Remove(DbPath)
		DbPath = dbPathReal
	}()
	DbPath = "/tmp"

	for _, test := range tests {
		err := test.errPrvdr()
		if err == nil {
			t.Errorf("Error is expected when %s is called\n", test.caseDesc)
		}

	}
}
func TestBucketInitialization(t *testing.T) {
	savedInit := Init
	dbPathReal := DbPath
	defer func() {
		os.Remove(DbPath)
		Init = savedInit
		DbPath = dbPathReal
	}()
	DbPath = "test_webhooks.db"
	expectedError := errors.New("weird error")
	Init = func(db *bbolt.DB, bucket string) error {
		return expectedError
	}
	for _, test := range tests {
		if test.initIsNotCalled {
			continue
		}
		err := test.errPrvdr()
		if err != expectedError {
			t.Errorf("Unexpected error for %s call, expected %v, got %v", test.caseDesc, expectedError, err)
		}

	}
}
