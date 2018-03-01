// Package models contains the types for schema ''.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
	"errors"
)

// CuptiActivityKindDriver represents a row from 'CUPTI_ACTIVITY_KIND_DRIVER'.
type CuptiActivityKindDriver struct {
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

// Exists determines if the CuptiActivityKindDriver exists in the database.
func (cakd *CuptiActivityKindDriver) Exists() bool {
	return cakd._exists
}

// Deleted provides information if the CuptiActivityKindDriver has been deleted from the database.
func (cakd *CuptiActivityKindDriver) Deleted() bool {
	return cakd._deleted
}

// Insert inserts the CuptiActivityKindDriver to the database.
func (cakd *CuptiActivityKindDriver) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if cakd._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO CUPTI_ACTIVITY_KIND_DRIVER (` +
		`cbid, start, end, processId, threadId, correlationId, returnValue` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, cakd.Cbid, cakd.Start, cakd.End, cakd.Processid, cakd.Threadid, cakd.Correlationid, cakd.Returnvalue)
	res, err := db.Exec(sqlstr, cakd.Cbid, cakd.Start, cakd.End, cakd.Processid, cakd.Threadid, cakd.Correlationid, cakd.Returnvalue)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	cakd.ID = sql.NullInt64(id)
	cakd._exists = true

	return nil
}

// Update updates the CuptiActivityKindDriver in the database.
func (cakd *CuptiActivityKindDriver) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !cakd._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if cakd._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE CUPTI_ACTIVITY_KIND_DRIVER SET ` +
		`cbid = ?, start = ?, end = ?, processId = ?, threadId = ?, correlationId = ?, returnValue = ?` +
		` WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, cakd.Cbid, cakd.Start, cakd.End, cakd.Processid, cakd.Threadid, cakd.Correlationid, cakd.Returnvalue, cakd.ID)
	_, err = db.Exec(sqlstr, cakd.Cbid, cakd.Start, cakd.End, cakd.Processid, cakd.Threadid, cakd.Correlationid, cakd.Returnvalue, cakd.ID)
	return err
}

// Save saves the CuptiActivityKindDriver to the database.
func (cakd *CuptiActivityKindDriver) Save(db XODB) error {
	if cakd.Exists() {
		return cakd.Update(db)
	}

	return cakd.Insert(db)
}

// Delete deletes the CuptiActivityKindDriver from the database.
func (cakd *CuptiActivityKindDriver) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !cakd._exists {
		return nil
	}

	// if deleted, bail
	if cakd._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM CUPTI_ACTIVITY_KIND_DRIVER WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, cakd.ID)
	_, err = db.Exec(sqlstr, cakd.ID)
	if err != nil {
		return err
	}

	// set deleted
	cakd._deleted = true

	return nil
}

// CuptiActivityKindDriverByID retrieves a row from 'CUPTI_ACTIVITY_KIND_DRIVER' as a CuptiActivityKindDriver.
//
// Generated from index 'CUPTI_ACTIVITY_KIND_DRIVER__id__pkey'.
func CuptiActivityKindDriverByID(db XODB, id sql.NullInt64) (*CuptiActivityKindDriver, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`_id_, cbid, start, end, processId, threadId, correlationId, returnValue ` +
		`FROM CUPTI_ACTIVITY_KIND_DRIVER ` +
		`WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, id)
	cakd := CuptiActivityKindDriver{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&cakd.ID, &cakd.Cbid, &cakd.Start, &cakd.End, &cakd.Processid, &cakd.Threadid, &cakd.Correlationid, &cakd.Returnvalue)
	if err != nil {
		return nil, err
	}

	return &cakd, nil
}

// CuptiActivityKindDriversByCorrelationid retrieves a row from 'CUPTI_ACTIVITY_KIND_DRIVER' as a CuptiActivityKindDriver.
//
// Generated from index 'INDEX_CUPTI_ACTIVITY_KIND_DRIVER'.
func CuptiActivityKindDriversByCorrelationid(db XODB, correlationid int) ([]*CuptiActivityKindDriver, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`_id_, cbid, start, end, processId, threadId, correlationId, returnValue ` +
		`FROM CUPTI_ACTIVITY_KIND_DRIVER ` +
		`WHERE correlationId = ?`

	// run query
	XOLog(sqlstr, correlationid)
	q, err := db.Query(sqlstr, correlationid)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*CuptiActivityKindDriver{}
	for q.Next() {
		cakd := CuptiActivityKindDriver{
			_exists: true,
		}

		// scan
		err = q.Scan(&cakd.ID, &cakd.Cbid, &cakd.Start, &cakd.End, &cakd.Processid, &cakd.Threadid, &cakd.Correlationid, &cakd.Returnvalue)
		if err != nil {
			return nil, err
		}

		res = append(res, &cakd)
	}

	return res, nil
}