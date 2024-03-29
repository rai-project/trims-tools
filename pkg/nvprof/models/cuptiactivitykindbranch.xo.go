// Package models contains the types for schema ''.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
	"errors"
)

// CuptiActivityKindBranch represents a row from 'CUPTI_ACTIVITY_KIND_BRANCH'.
type CuptiActivityKindBranch struct {
	ID              sql.NullInt64 `json:"_id_"`            // _id_
	Sourcelocatorid int           `json:"sourceLocatorId"` // sourceLocatorId
	Correlationid   int           `json:"correlationId"`   // correlationId
	Functionid      int           `json:"functionId"`      // functionId
	Pcoffset        int           `json:"pcOffset"`        // pcOffset
	Diverged        int           `json:"diverged"`        // diverged
	Threadsexecuted int           `json:"threadsExecuted"` // threadsExecuted
	Executed        int           `json:"executed"`        // executed

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the CuptiActivityKindBranch exists in the database.
func (cakb *CuptiActivityKindBranch) Exists() bool {
	return cakb._exists
}

// Deleted provides information if the CuptiActivityKindBranch has been deleted from the database.
func (cakb *CuptiActivityKindBranch) Deleted() bool {
	return cakb._deleted
}

// Insert inserts the CuptiActivityKindBranch to the database.
func (cakb *CuptiActivityKindBranch) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if cakb._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO CUPTI_ACTIVITY_KIND_BRANCH (` +
		`sourceLocatorId, correlationId, functionId, pcOffset, diverged, threadsExecuted, executed` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, cakb.Sourcelocatorid, cakb.Correlationid, cakb.Functionid, cakb.Pcoffset, cakb.Diverged, cakb.Threadsexecuted, cakb.Executed)
	res, err := db.Exec(sqlstr, cakb.Sourcelocatorid, cakb.Correlationid, cakb.Functionid, cakb.Pcoffset, cakb.Diverged, cakb.Threadsexecuted, cakb.Executed)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	cakb.ID = sql.NullInt64(id)
	cakb._exists = true

	return nil
}

// Update updates the CuptiActivityKindBranch in the database.
func (cakb *CuptiActivityKindBranch) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !cakb._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if cakb._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE CUPTI_ACTIVITY_KIND_BRANCH SET ` +
		`sourceLocatorId = ?, correlationId = ?, functionId = ?, pcOffset = ?, diverged = ?, threadsExecuted = ?, executed = ?` +
		` WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, cakb.Sourcelocatorid, cakb.Correlationid, cakb.Functionid, cakb.Pcoffset, cakb.Diverged, cakb.Threadsexecuted, cakb.Executed, cakb.ID)
	_, err = db.Exec(sqlstr, cakb.Sourcelocatorid, cakb.Correlationid, cakb.Functionid, cakb.Pcoffset, cakb.Diverged, cakb.Threadsexecuted, cakb.Executed, cakb.ID)
	return err
}

// Save saves the CuptiActivityKindBranch to the database.
func (cakb *CuptiActivityKindBranch) Save(db XODB) error {
	if cakb.Exists() {
		return cakb.Update(db)
	}

	return cakb.Insert(db)
}

// Delete deletes the CuptiActivityKindBranch from the database.
func (cakb *CuptiActivityKindBranch) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !cakb._exists {
		return nil
	}

	// if deleted, bail
	if cakb._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM CUPTI_ACTIVITY_KIND_BRANCH WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, cakb.ID)
	_, err = db.Exec(sqlstr, cakb.ID)
	if err != nil {
		return err
	}

	// set deleted
	cakb._deleted = true

	return nil
}

// CuptiActivityKindBranchByID retrieves a row from 'CUPTI_ACTIVITY_KIND_BRANCH' as a CuptiActivityKindBranch.
//
// Generated from index 'CUPTI_ACTIVITY_KIND_BRANCH__id__pkey'.
func CuptiActivityKindBranchByID(db XODB, id sql.NullInt64) (*CuptiActivityKindBranch, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`_id_, sourceLocatorId, correlationId, functionId, pcOffset, diverged, threadsExecuted, executed ` +
		`FROM CUPTI_ACTIVITY_KIND_BRANCH ` +
		`WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, id)
	cakb := CuptiActivityKindBranch{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&cakb.ID, &cakb.Sourcelocatorid, &cakb.Correlationid, &cakb.Functionid, &cakb.Pcoffset, &cakb.Diverged, &cakb.Threadsexecuted, &cakb.Executed)
	if err != nil {
		return nil, err
	}

	return &cakb, nil
}

// CuptiActivityKindBranchesByCorrelationid retrieves a row from 'CUPTI_ACTIVITY_KIND_BRANCH' as a CuptiActivityKindBranch.
//
// Generated from index 'INDEX_CUPTI_ACTIVITY_KIND_BRANCH'.
func CuptiActivityKindBranchesByCorrelationid(db XODB, correlationid int) ([]*CuptiActivityKindBranch, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`_id_, sourceLocatorId, correlationId, functionId, pcOffset, diverged, threadsExecuted, executed ` +
		`FROM CUPTI_ACTIVITY_KIND_BRANCH ` +
		`WHERE correlationId = ?`

	// run query
	XOLog(sqlstr, correlationid)
	q, err := db.Query(sqlstr, correlationid)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*CuptiActivityKindBranch{}
	for q.Next() {
		cakb := CuptiActivityKindBranch{
			_exists: true,
		}

		// scan
		err = q.Scan(&cakb.ID, &cakb.Sourcelocatorid, &cakb.Correlationid, &cakb.Functionid, &cakb.Pcoffset, &cakb.Diverged, &cakb.Threadsexecuted, &cakb.Executed)
		if err != nil {
			return nil, err
		}

		res = append(res, &cakb)
	}

	return res, nil
}
