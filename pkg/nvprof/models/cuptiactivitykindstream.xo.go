// Package models contains the types for schema ''.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
	"errors"
)

// CuptiActivityKindStream represents a row from 'CUPTI_ACTIVITY_KIND_STREAM'.
type CuptiActivityKindStream struct {
	ID            sql.NullInt64 `json:"_id_"`          // _id_
	Contextid     int           `json:"contextId"`     // contextId
	Streamid      int           `json:"streamId"`      // streamId
	Priority      int           `json:"priority"`      // priority
	Flag          int           `json:"flag"`          // flag
	Correlationid int           `json:"correlationId"` // correlationId

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the CuptiActivityKindStream exists in the database.
func (caks *CuptiActivityKindStream) Exists() bool {
	return caks._exists
}

// Deleted provides information if the CuptiActivityKindStream has been deleted from the database.
func (caks *CuptiActivityKindStream) Deleted() bool {
	return caks._deleted
}

// Insert inserts the CuptiActivityKindStream to the database.
func (caks *CuptiActivityKindStream) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if caks._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO CUPTI_ACTIVITY_KIND_STREAM (` +
		`contextId, streamId, priority, flag, correlationId` +
		`) VALUES (` +
		`?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, caks.Contextid, caks.Streamid, caks.Priority, caks.Flag, caks.Correlationid)
	res, err := db.Exec(sqlstr, caks.Contextid, caks.Streamid, caks.Priority, caks.Flag, caks.Correlationid)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	caks.ID = sql.NullInt64(id)
	caks._exists = true

	return nil
}

// Update updates the CuptiActivityKindStream in the database.
func (caks *CuptiActivityKindStream) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !caks._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if caks._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE CUPTI_ACTIVITY_KIND_STREAM SET ` +
		`contextId = ?, streamId = ?, priority = ?, flag = ?, correlationId = ?` +
		` WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, caks.Contextid, caks.Streamid, caks.Priority, caks.Flag, caks.Correlationid, caks.ID)
	_, err = db.Exec(sqlstr, caks.Contextid, caks.Streamid, caks.Priority, caks.Flag, caks.Correlationid, caks.ID)
	return err
}

// Save saves the CuptiActivityKindStream to the database.
func (caks *CuptiActivityKindStream) Save(db XODB) error {
	if caks.Exists() {
		return caks.Update(db)
	}

	return caks.Insert(db)
}

// Delete deletes the CuptiActivityKindStream from the database.
func (caks *CuptiActivityKindStream) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !caks._exists {
		return nil
	}

	// if deleted, bail
	if caks._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM CUPTI_ACTIVITY_KIND_STREAM WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, caks.ID)
	_, err = db.Exec(sqlstr, caks.ID)
	if err != nil {
		return err
	}

	// set deleted
	caks._deleted = true

	return nil
}

// CuptiActivityKindStreamByID retrieves a row from 'CUPTI_ACTIVITY_KIND_STREAM' as a CuptiActivityKindStream.
//
// Generated from index 'CUPTI_ACTIVITY_KIND_STREAM__id__pkey'.
func CuptiActivityKindStreamByID(db XODB, id sql.NullInt64) (*CuptiActivityKindStream, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`_id_, contextId, streamId, priority, flag, correlationId ` +
		`FROM CUPTI_ACTIVITY_KIND_STREAM ` +
		`WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, id)
	caks := CuptiActivityKindStream{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&caks.ID, &caks.Contextid, &caks.Streamid, &caks.Priority, &caks.Flag, &caks.Correlationid)
	if err != nil {
		return nil, err
	}

	return &caks, nil
}

// CuptiActivityKindStreamsByCorrelationid retrieves a row from 'CUPTI_ACTIVITY_KIND_STREAM' as a CuptiActivityKindStream.
//
// Generated from index 'INDEX_CUPTI_ACTIVITY_KIND_STREAM'.
func CuptiActivityKindStreamsByCorrelationid(db XODB, correlationid int) ([]*CuptiActivityKindStream, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`_id_, contextId, streamId, priority, flag, correlationId ` +
		`FROM CUPTI_ACTIVITY_KIND_STREAM ` +
		`WHERE correlationId = ?`

	// run query
	XOLog(sqlstr, correlationid)
	q, err := db.Query(sqlstr, correlationid)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*CuptiActivityKindStream{}
	for q.Next() {
		caks := CuptiActivityKindStream{
			_exists: true,
		}

		// scan
		err = q.Scan(&caks.ID, &caks.Contextid, &caks.Streamid, &caks.Priority, &caks.Flag, &caks.Correlationid)
		if err != nil {
			return nil, err
		}

		res = append(res, &caks)
	}

	return res, nil
}
