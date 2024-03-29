// Package models contains the types for schema ''.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
	"errors"
)

// CuptiActivityKindInstructionExecution represents a row from 'CUPTI_ACTIVITY_KIND_INSTRUCTION_EXECUTION'.
type CuptiActivityKindInstructionExecution struct {
	ID                        sql.NullInt64 `json:"_id_"`                      // _id_
	Flags                     int           `json:"flags"`                     // flags
	Sourcelocatorid           int           `json:"sourceLocatorId"`           // sourceLocatorId
	Correlationid             int           `json:"correlationId"`             // correlationId
	Functionid                int           `json:"functionId"`                // functionId
	Pcoffset                  int           `json:"pcOffset"`                  // pcOffset
	Threadsexecuted           int           `json:"threadsExecuted"`           // threadsExecuted
	Notpredoffthreadsexecuted int           `json:"notPredOffThreadsExecuted"` // notPredOffThreadsExecuted
	Executed                  int           `json:"executed"`                  // executed

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the CuptiActivityKindInstructionExecution exists in the database.
func (cakie *CuptiActivityKindInstructionExecution) Exists() bool {
	return cakie._exists
}

// Deleted provides information if the CuptiActivityKindInstructionExecution has been deleted from the database.
func (cakie *CuptiActivityKindInstructionExecution) Deleted() bool {
	return cakie._deleted
}

// Insert inserts the CuptiActivityKindInstructionExecution to the database.
func (cakie *CuptiActivityKindInstructionExecution) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if cakie._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO CUPTI_ACTIVITY_KIND_INSTRUCTION_EXECUTION (` +
		`flags, sourceLocatorId, correlationId, functionId, pcOffset, threadsExecuted, notPredOffThreadsExecuted, executed` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, cakie.Flags, cakie.Sourcelocatorid, cakie.Correlationid, cakie.Functionid, cakie.Pcoffset, cakie.Threadsexecuted, cakie.Notpredoffthreadsexecuted, cakie.Executed)
	res, err := db.Exec(sqlstr, cakie.Flags, cakie.Sourcelocatorid, cakie.Correlationid, cakie.Functionid, cakie.Pcoffset, cakie.Threadsexecuted, cakie.Notpredoffthreadsexecuted, cakie.Executed)
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

// Update updates the CuptiActivityKindInstructionExecution in the database.
func (cakie *CuptiActivityKindInstructionExecution) Update(db XODB) error {
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
	const sqlstr = `UPDATE CUPTI_ACTIVITY_KIND_INSTRUCTION_EXECUTION SET ` +
		`flags = ?, sourceLocatorId = ?, correlationId = ?, functionId = ?, pcOffset = ?, threadsExecuted = ?, notPredOffThreadsExecuted = ?, executed = ?` +
		` WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, cakie.Flags, cakie.Sourcelocatorid, cakie.Correlationid, cakie.Functionid, cakie.Pcoffset, cakie.Threadsexecuted, cakie.Notpredoffthreadsexecuted, cakie.Executed, cakie.ID)
	_, err = db.Exec(sqlstr, cakie.Flags, cakie.Sourcelocatorid, cakie.Correlationid, cakie.Functionid, cakie.Pcoffset, cakie.Threadsexecuted, cakie.Notpredoffthreadsexecuted, cakie.Executed, cakie.ID)
	return err
}

// Save saves the CuptiActivityKindInstructionExecution to the database.
func (cakie *CuptiActivityKindInstructionExecution) Save(db XODB) error {
	if cakie.Exists() {
		return cakie.Update(db)
	}

	return cakie.Insert(db)
}

// Delete deletes the CuptiActivityKindInstructionExecution from the database.
func (cakie *CuptiActivityKindInstructionExecution) Delete(db XODB) error {
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
	const sqlstr = `DELETE FROM CUPTI_ACTIVITY_KIND_INSTRUCTION_EXECUTION WHERE _id_ = ?`

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

// CuptiActivityKindInstructionExecutionByID retrieves a row from 'CUPTI_ACTIVITY_KIND_INSTRUCTION_EXECUTION' as a CuptiActivityKindInstructionExecution.
//
// Generated from index 'CUPTI_ACTIVITY_KIND_INSTRUCTION_EXECUTION__id__pkey'.
func CuptiActivityKindInstructionExecutionByID(db XODB, id sql.NullInt64) (*CuptiActivityKindInstructionExecution, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`_id_, flags, sourceLocatorId, correlationId, functionId, pcOffset, threadsExecuted, notPredOffThreadsExecuted, executed ` +
		`FROM CUPTI_ACTIVITY_KIND_INSTRUCTION_EXECUTION ` +
		`WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, id)
	cakie := CuptiActivityKindInstructionExecution{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&cakie.ID, &cakie.Flags, &cakie.Sourcelocatorid, &cakie.Correlationid, &cakie.Functionid, &cakie.Pcoffset, &cakie.Threadsexecuted, &cakie.Notpredoffthreadsexecuted, &cakie.Executed)
	if err != nil {
		return nil, err
	}

	return &cakie, nil
}

// CuptiActivityKindInstructionExecutionsByCorrelationid retrieves a row from 'CUPTI_ACTIVITY_KIND_INSTRUCTION_EXECUTION' as a CuptiActivityKindInstructionExecution.
//
// Generated from index 'INDEX_CUPTI_ACTIVITY_KIND_INSTRUCTION_EXECUTION'.
func CuptiActivityKindInstructionExecutionsByCorrelationid(db XODB, correlationid int) ([]*CuptiActivityKindInstructionExecution, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`_id_, flags, sourceLocatorId, correlationId, functionId, pcOffset, threadsExecuted, notPredOffThreadsExecuted, executed ` +
		`FROM CUPTI_ACTIVITY_KIND_INSTRUCTION_EXECUTION ` +
		`WHERE correlationId = ?`

	// run query
	XOLog(sqlstr, correlationid)
	q, err := db.Query(sqlstr, correlationid)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*CuptiActivityKindInstructionExecution{}
	for q.Next() {
		cakie := CuptiActivityKindInstructionExecution{
			_exists: true,
		}

		// scan
		err = q.Scan(&cakie.ID, &cakie.Flags, &cakie.Sourcelocatorid, &cakie.Correlationid, &cakie.Functionid, &cakie.Pcoffset, &cakie.Threadsexecuted, &cakie.Notpredoffthreadsexecuted, &cakie.Executed)
		if err != nil {
			return nil, err
		}

		res = append(res, &cakie)
	}

	return res, nil
}
