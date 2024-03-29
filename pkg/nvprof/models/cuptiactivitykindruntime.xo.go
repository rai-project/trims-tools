// Package models contains the types for schema ''.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
	"errors"
)

// CuptiActivityKindRuntime represents a row from 'CUPTI_ACTIVITY_KIND_RUNTIME'.
type CuptiActivityKindRuntime struct {
	ID            sql.NullInt64 `json:"_id_"`          // _id_
	Cbid          int           `json:"cbid"`          // cbid
	Start         int           `json:"start"`         // start
	End           int           `json:"end"`           // end
	Processid     int           `json:"processId"`     // processId
	Threadid      int           `json:"threadId"`      // threadId
	Correlationid int           `json:"correlationId"` // correlationId
	Returnvalue   int           `json:"returnValue"`   // returnValue

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the CuptiActivityKindRuntime exists in the database.
func (cakr *CuptiActivityKindRuntime) Exists() bool {
	return cakr._exists
}

// Deleted provides information if the CuptiActivityKindRuntime has been deleted from the database.
func (cakr *CuptiActivityKindRuntime) Deleted() bool {
	return cakr._deleted
}

// Insert inserts the CuptiActivityKindRuntime to the database.
func (cakr *CuptiActivityKindRuntime) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if cakr._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO CUPTI_ACTIVITY_KIND_RUNTIME (` +
		`cbid, start, end, processId, threadId, correlationId, returnValue` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, cakr.Cbid, cakr.Start, cakr.End, cakr.Processid, cakr.Threadid, cakr.Correlationid, cakr.Returnvalue)
	res, err := db.Exec(sqlstr, cakr.Cbid, cakr.Start, cakr.End, cakr.Processid, cakr.Threadid, cakr.Correlationid, cakr.Returnvalue)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	cakr.ID = sql.NullInt64(id)
	cakr._exists = true

	return nil
}

// Update updates the CuptiActivityKindRuntime in the database.
func (cakr *CuptiActivityKindRuntime) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !cakr._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if cakr._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE CUPTI_ACTIVITY_KIND_RUNTIME SET ` +
		`cbid = ?, start = ?, end = ?, processId = ?, threadId = ?, correlationId = ?, returnValue = ?` +
		` WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, cakr.Cbid, cakr.Start, cakr.End, cakr.Processid, cakr.Threadid, cakr.Correlationid, cakr.Returnvalue, cakr.ID)
	_, err = db.Exec(sqlstr, cakr.Cbid, cakr.Start, cakr.End, cakr.Processid, cakr.Threadid, cakr.Correlationid, cakr.Returnvalue, cakr.ID)
	return err
}

// Save saves the CuptiActivityKindRuntime to the database.
func (cakr *CuptiActivityKindRuntime) Save(db XODB) error {
	if cakr.Exists() {
		return cakr.Update(db)
	}

	return cakr.Insert(db)
}

// Delete deletes the CuptiActivityKindRuntime from the database.
func (cakr *CuptiActivityKindRuntime) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !cakr._exists {
		return nil
	}

	// if deleted, bail
	if cakr._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM CUPTI_ACTIVITY_KIND_RUNTIME WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, cakr.ID)
	_, err = db.Exec(sqlstr, cakr.ID)
	if err != nil {
		return err
	}

	// set deleted
	cakr._deleted = true

	return nil
}

// CuptiActivityKindRuntimeByID retrieves a row from 'CUPTI_ACTIVITY_KIND_RUNTIME' as a CuptiActivityKindRuntime.
//
// Generated from index 'CUPTI_ACTIVITY_KIND_RUNTIME__id__pkey'.
func CuptiActivityKindRuntimeByID(db XODB, id sql.NullInt64) (*CuptiActivityKindRuntime, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`_id_, cbid, start, end, processId, threadId, correlationId, returnValue ` +
		`FROM CUPTI_ACTIVITY_KIND_RUNTIME ` +
		`WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, id)
	cakr := CuptiActivityKindRuntime{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&cakr.ID, &cakr.Cbid, &cakr.Start, &cakr.End, &cakr.Processid, &cakr.Threadid, &cakr.Correlationid, &cakr.Returnvalue)
	if err != nil {
		return nil, err
	}

	return &cakr, nil
}

// CuptiActivityKindRuntimesByCorrelationid retrieves a row from 'CUPTI_ACTIVITY_KIND_RUNTIME' as a CuptiActivityKindRuntime.
//
// Generated from index 'INDEX_CUPTI_ACTIVITY_KIND_RUNTIME'.
func CuptiActivityKindRuntimesByCorrelationid(db XODB, correlationid int) ([]*CuptiActivityKindRuntime, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`_id_, cbid, start, end, processId, threadId, correlationId, returnValue ` +
		`FROM CUPTI_ACTIVITY_KIND_RUNTIME ` +
		`WHERE correlationId = ?`

	// run query
	XOLog(sqlstr, correlationid)
	q, err := db.Query(sqlstr, correlationid)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*CuptiActivityKindRuntime{}
	for q.Next() {
		cakr := CuptiActivityKindRuntime{
			_exists: true,
		}

		// scan
		err = q.Scan(&cakr.ID, &cakr.Cbid, &cakr.Start, &cakr.End, &cakr.Processid, &cakr.Threadid, &cakr.Correlationid, &cakr.Returnvalue)
		if err != nil {
			return nil, err
		}

		res = append(res, &cakr)
	}

	return res, nil
}
