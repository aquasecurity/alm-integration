package dbservice

import (
	"encoding/json"

	bolt "go.etcd.io/bbolt"
)

func AggregateScans(output string,
	currentScan map[string]string,
	scansPerTicket int,
	ignoreTheQuantity bool) ([]map[string]string, error) {
	mutex.Lock()
	defer mutex.Unlock()

	db, err := bolt.Open(DbPath, 0666, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = Init(db, dbBucketAggregator)
	if err != nil {
		return nil, err
	}

	aggregatedScans := make([]map[string]string, 0, scansPerTicket)
	if len(currentScan) > 0 {
		aggregatedScans = append(aggregatedScans, currentScan)
	}
	currentValue, err := dbSelect(db, dbBucketAggregator, output)
	if err != nil {
		return nil, err
	}

	if len(currentValue) > 0 {
		var savedScans []map[string]string
		err = json.Unmarshal(currentValue, &savedScans)
		if err != nil {
			return nil, err
		}
		aggregatedScans = append(aggregatedScans, savedScans...)
	}

	if ignoreTheQuantity || len(aggregatedScans) < scansPerTicket {
		saving, err := json.Marshal(aggregatedScans)
		if err != nil {
			return nil, err
		}

		err = dbInsert(db, dbBucketAggregator, []byte(output), saving)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	dbInsert(db, dbBucketAggregator, []byte(output), nil)
	return aggregatedScans, nil
}
