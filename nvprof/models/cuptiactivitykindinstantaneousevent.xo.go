// Package models contains the types for schema ''.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
	"errors"
)

// CuptiActivityKindInstantaneousEvent represents a row from 'CUPTI_ACTIVITY_KIND_INSTANTANEOUS_EVENT'.
type CuptiActivityKindInstantaneousEvent struct {
	ID        sql.NullInt64 `json:"_id_"`      // _id_
	ID        int           `json:"id"`        // id
	Value     int           `json:"value"`     // value
	Timestamp int           `json:"timestamp"` // timestamp
	Deviceid  int           `json:"deviceId"`  // deviceId

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the CuptiActivityKindInstantaneousEvent exists in the database.
func (cakie *CuptiActivityKindInstantaneousEvent) Exists() bool {
	return cakie._exists
}

// Deleted provides information if the CuptiActivityKindInstantaneousEvent has been deleted from the database.
func (cakie *CuptiActivityKindInstantaneousEvent) Deleted() bool {
	return cakie._deleted
}

// Insert inserts the CuptiActivityKindInstantaneousEvent to the database.
func (cakie *CuptiActivityKindInstantaneousEvent) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if cakie._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO CUPTI_ACTIVITY_KIND_INSTANTANEOUS_EVENT (` +
		`value, timestamp, deviceId` +
		`) VALUES (` +
		`?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, cakie.Value, cakie.Timestamp, cakie.Deviceid)
	res, err := db.Exec(sqlstr, cakie.Value, cakie.Timestamp, cakie.Deviceid)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	cakie.ID = sql.NullInt64(id)
	cakie._exists = true

	return nil
}

// Update updates the CuptiActivityKindInstantaneousEvent in the database.
func (cakie *CuptiActivityKindInstantaneousEvent) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !cakie._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if cakie._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE CUPTI_ACTIVITY_KIND_INSTANTANEOUS_EVENT SET ` +
		`value = ?, timestamp = ?, deviceId = ?` +
		` WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, cakie.Value, cakie.Timestamp, cakie.Deviceid, cakie.ID)
	_, err = db.Exec(sqlstr, cakie.Value, cakie.Timestamp, cakie.Deviceid, cakie.ID)
	return err
}

// Save saves the CuptiActivityKindInstantaneousEvent to the database.
func (cakie *CuptiActivityKindInstantaneousEvent) Save(db XODB) error {
	if cakie.Exists() {
		return cakie.Update(db)
	}

	return cakie.Insert(db)
}

// Delete deletes the CuptiActivityKindInstantaneousEvent from the database.
func (cakie *CuptiActivityKindInstantaneousEvent) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !cakie._exists {
		return nil
	}

	// if deleted, bail
	if cakie._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM CUPTI_ACTIVITY_KIND_INSTANTANEOUS_EVENT WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, cakie.ID)
	_, err = db.Exec(sqlstr, cakie.ID)
	if err != nil {
		return err
	}

	// set deleted
	cakie._deleted = true

	return nil
}

// CuptiActivityKindInstantaneousEventByID retrieves a row from 'CUPTI_ACTIVITY_KIND_INSTANTANEOUS_EVENT' as a CuptiActivityKindInstantaneousEvent.
//
// Generated from index 'CUPTI_ACTIVITY_KIND_INSTANTANEOUS_EVENT__id__pkey'.
func CuptiActivityKindInstantaneousEventByID(db XODB, id sql.NullInt64) (*CuptiActivityKindInstantaneousEvent, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`_id_, id, value, timestamp, deviceId ` +
		`FROM CUPTI_ACTIVITY_KIND_INSTANTANEOUS_EVENT ` +
		`WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, id)
	cakie := CuptiActivityKindInstantaneousEvent{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&cakie.ID, &cakie.ID, &cakie.Value, &cakie.Timestamp, &cakie.Deviceid)
	if err != nil {
		return nil, err
	}

	return &cakie, nil
}

// CuptiActivityKindInstantaneousEventsByTimestamp retrieves a row from 'CUPTI_ACTIVITY_KIND_INSTANTANEOUS_EVENT' as a CuptiActivityKindInstantaneousEvent.
//
// Generated from index 'INDEX_CUPTI_ACTIVITY_KIND_INSTANTANEOUS_EVENT'.
func CuptiActivityKindInstantaneousEventsByTimestamp(db XODB, timestamp int) ([]*CuptiActivityKindInstantaneousEvent, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`_id_, id, value, timestamp, deviceId ` +
		`FROM CUPTI_ACTIVITY_KIND_INSTANTANEOUS_EVENT ` +
		`WHERE timestamp = ?`

	// run query
	XOLog(sqlstr, timestamp)
	q, err := db.Query(sqlstr, timestamp)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*CuptiActivityKindInstantaneousEvent{}
	for q.Next() {
		cakie := CuptiActivityKindInstantaneousEvent{
			_exists: true,
		}

		// scan
		err = q.Scan(&cakie.ID, &cakie.ID, &cakie.Value, &cakie.Timestamp, &cakie.Deviceid)
		if err != nil {
			return nil, err
		}

		res = append(res, &cakie)
	}

	return res, nil
}
