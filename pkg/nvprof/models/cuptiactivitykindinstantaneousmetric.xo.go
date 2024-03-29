// Package models contains the types for schema ''.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
	"errors"
)

// CuptiActivityKindInstantaneousMetric represents a row from 'CUPTI_ACTIVITY_KIND_INSTANTANEOUS_METRIC'.
type CuptiActivityKindInstantaneousMetric struct {
	ID        sql.NullInt64 `json:"_id_"`      // _id_
	ID        int           `json:"id"`        // id
	Value     []byte        `json:"value"`     // value
	Timestamp int           `json:"timestamp"` // timestamp
	Deviceid  int           `json:"deviceId"`  // deviceId
	Flags     int           `json:"flags"`     // flags

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the CuptiActivityKindInstantaneousMetric exists in the database.
func (cakim *CuptiActivityKindInstantaneousMetric) Exists() bool {
	return cakim._exists
}

// Deleted provides information if the CuptiActivityKindInstantaneousMetric has been deleted from the database.
func (cakim *CuptiActivityKindInstantaneousMetric) Deleted() bool {
	return cakim._deleted
}

// Insert inserts the CuptiActivityKindInstantaneousMetric to the database.
func (cakim *CuptiActivityKindInstantaneousMetric) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if cakim._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO CUPTI_ACTIVITY_KIND_INSTANTANEOUS_METRIC (` +
		`value, timestamp, deviceId, flags` +
		`) VALUES (` +
		`?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, cakim.Value, cakim.Timestamp, cakim.Deviceid, cakim.Flags)
	res, err := db.Exec(sqlstr, cakim.Value, cakim.Timestamp, cakim.Deviceid, cakim.Flags)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	cakim.ID = sql.NullInt64(id)
	cakim._exists = true

	return nil
}

// Update updates the CuptiActivityKindInstantaneousMetric in the database.
func (cakim *CuptiActivityKindInstantaneousMetric) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !cakim._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if cakim._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE CUPTI_ACTIVITY_KIND_INSTANTANEOUS_METRIC SET ` +
		`value = ?, timestamp = ?, deviceId = ?, flags = ?` +
		` WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, cakim.Value, cakim.Timestamp, cakim.Deviceid, cakim.Flags, cakim.ID)
	_, err = db.Exec(sqlstr, cakim.Value, cakim.Timestamp, cakim.Deviceid, cakim.Flags, cakim.ID)
	return err
}

// Save saves the CuptiActivityKindInstantaneousMetric to the database.
func (cakim *CuptiActivityKindInstantaneousMetric) Save(db XODB) error {
	if cakim.Exists() {
		return cakim.Update(db)
	}

	return cakim.Insert(db)
}

// Delete deletes the CuptiActivityKindInstantaneousMetric from the database.
func (cakim *CuptiActivityKindInstantaneousMetric) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !cakim._exists {
		return nil
	}

	// if deleted, bail
	if cakim._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM CUPTI_ACTIVITY_KIND_INSTANTANEOUS_METRIC WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, cakim.ID)
	_, err = db.Exec(sqlstr, cakim.ID)
	if err != nil {
		return err
	}

	// set deleted
	cakim._deleted = true

	return nil
}

// CuptiActivityKindInstantaneousMetricByID retrieves a row from 'CUPTI_ACTIVITY_KIND_INSTANTANEOUS_METRIC' as a CuptiActivityKindInstantaneousMetric.
//
// Generated from index 'CUPTI_ACTIVITY_KIND_INSTANTANEOUS_METRIC__id__pkey'.
func CuptiActivityKindInstantaneousMetricByID(db XODB, id sql.NullInt64) (*CuptiActivityKindInstantaneousMetric, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`_id_, id, value, timestamp, deviceId, flags ` +
		`FROM CUPTI_ACTIVITY_KIND_INSTANTANEOUS_METRIC ` +
		`WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, id)
	cakim := CuptiActivityKindInstantaneousMetric{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&cakim.ID, &cakim.ID, &cakim.Value, &cakim.Timestamp, &cakim.Deviceid, &cakim.Flags)
	if err != nil {
		return nil, err
	}

	return &cakim, nil
}

// CuptiActivityKindInstantaneousMetricsByTimestamp retrieves a row from 'CUPTI_ACTIVITY_KIND_INSTANTANEOUS_METRIC' as a CuptiActivityKindInstantaneousMetric.
//
// Generated from index 'INDEX_CUPTI_ACTIVITY_KIND_INSTANTANEOUS_METRIC'.
func CuptiActivityKindInstantaneousMetricsByTimestamp(db XODB, timestamp int) ([]*CuptiActivityKindInstantaneousMetric, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`_id_, id, value, timestamp, deviceId, flags ` +
		`FROM CUPTI_ACTIVITY_KIND_INSTANTANEOUS_METRIC ` +
		`WHERE timestamp = ?`

	// run query
	XOLog(sqlstr, timestamp)
	q, err := db.Query(sqlstr, timestamp)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*CuptiActivityKindInstantaneousMetric{}
	for q.Next() {
		cakim := CuptiActivityKindInstantaneousMetric{
			_exists: true,
		}

		// scan
		err = q.Scan(&cakim.ID, &cakim.ID, &cakim.Value, &cakim.Timestamp, &cakim.Deviceid, &cakim.Flags)
		if err != nil {
			return nil, err
		}

		res = append(res, &cakim)
	}

	return res, nil
}
