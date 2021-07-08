package dbservice

import "testing"

func TestChangeDbPath(t *testing.T) {
	testPath := "/tmp/test.db"
	storedPath := DbPath
	ChangeDbPath(testPath)
	defer func() {
		ChangeDbPath(storedPath)
	}()
	if DbPath != testPath {
		t.Errorf("path is not configured correctly, expected: %s, got %s", testPath, DbPath)
	}
}
