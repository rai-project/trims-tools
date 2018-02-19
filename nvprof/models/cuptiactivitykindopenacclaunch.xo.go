// Package models contains the types for schema ''.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
	"errors"
)

// CuptiActivityKindOpenaccLaunch represents a row from 'CUPTI_ACTIVITY_KIND_OPENACC_LAUNCH'.
type CuptiActivityKindOpenaccLaunch struct {
	ID              sql.NullInt64 `json:"_id_"`            // _id_
	Eventkind       int           `json:"eventKind"`       // eventKind
	Parentconstruct int           `json:"parentConstruct"` // parentConstruct
	Version         int           `json:"version"`         // version
	Implicit        int           `json:"implicit"`        // implicit
	Devicetype      int           `json:"deviceType"`      // deviceType
	Devicenumber    int           `json:"deviceNumber"`    // deviceNumber
	Threadid        int           `json:"threadId"`        // threadId
	Async           int           `json:"async"`           // async
	Asyncmap        int           `json:"asyncMap"`        // asyncMap
	Lineno          int           `json:"lineNo"`          // lineNo
	Endlineno       int           `json:"endLineNo"`       // endLineNo
	Funclineno      int           `json:"funcLineNo"`      // funcLineNo
	Funcendlineno   int           `json:"funcEndLineNo"`   // funcEndLineNo
	Start           int           `json:"start"`           // start
	End             int           `json:"end"`             // end
	Cudeviceid      int           `json:"cuDeviceId"`      // cuDeviceId
	Cucontextid     int           `json:"cuContextId"`     // cuContextId
	Custreamid      int           `json:"cuStreamId"`      // cuStreamId
	Cuprocessid     int           `json:"cuProcessId"`     // cuProcessId
	Cuthreadid      int           `json:"cuThreadId"`      // cuThreadId
	Externalid      int           `json:"externalId"`      // externalId
	Srcfile         int           `json:"srcFile"`         // srcFile
	Funcname        int           `json:"funcName"`        // funcName
	Numgangs        int           `json:"numGangs"`        // numGangs
	Numworkers      int           `json:"numWorkers"`      // numWorkers
	Vectorlength    int           `json:"vectorLength"`    // vectorLength
	Kernelname      int           `json:"kernelName"`      // kernelName

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the CuptiActivityKindOpenaccLaunch exists in the database.
func (cakol *CuptiActivityKindOpenaccLaunch) Exists() bool {
	return cakol._exists
}

// Deleted provides information if the CuptiActivityKindOpenaccLaunch has been deleted from the database.
func (cakol *CuptiActivityKindOpenaccLaunch) Deleted() bool {
	return cakol._deleted
}

// Insert inserts the CuptiActivityKindOpenaccLaunch to the database.
func (cakol *CuptiActivityKindOpenaccLaunch) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if cakol._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO CUPTI_ACTIVITY_KIND_OPENACC_LAUNCH (` +
		`eventKind, parentConstruct, version, implicit, deviceType, deviceNumber, threadId, async, asyncMap, lineNo, endLineNo, funcLineNo, funcEndLineNo, start, end, cuDeviceId, cuContextId, cuStreamId, cuProcessId, cuThreadId, externalId, srcFile, funcName, numGangs, numWorkers, vectorLength, kernelName` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, cakol.Eventkind, cakol.Parentconstruct, cakol.Version, cakol.Implicit, cakol.Devicetype, cakol.Devicenumber, cakol.Threadid, cakol.Async, cakol.Asyncmap, cakol.Lineno, cakol.Endlineno, cakol.Funclineno, cakol.Funcendlineno, cakol.Start, cakol.End, cakol.Cudeviceid, cakol.Cucontextid, cakol.Custreamid, cakol.Cuprocessid, cakol.Cuthreadid, cakol.Externalid, cakol.Srcfile, cakol.Funcname, cakol.Numgangs, cakol.Numworkers, cakol.Vectorlength, cakol.Kernelname)
	res, err := db.Exec(sqlstr, cakol.Eventkind, cakol.Parentconstruct, cakol.Version, cakol.Implicit, cakol.Devicetype, cakol.Devicenumber, cakol.Threadid, cakol.Async, cakol.Asyncmap, cakol.Lineno, cakol.Endlineno, cakol.Funclineno, cakol.Funcendlineno, cakol.Start, cakol.End, cakol.Cudeviceid, cakol.Cucontextid, cakol.Custreamid, cakol.Cuprocessid, cakol.Cuthreadid, cakol.Externalid, cakol.Srcfile, cakol.Funcname, cakol.Numgangs, cakol.Numworkers, cakol.Vectorlength, cakol.Kernelname)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	cakol.ID = sql.NullInt64(id)
	cakol._exists = true

	return nil
}

// Update updates the CuptiActivityKindOpenaccLaunch in the database.
func (cakol *CuptiActivityKindOpenaccLaunch) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !cakol._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if cakol._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE CUPTI_ACTIVITY_KIND_OPENACC_LAUNCH SET ` +
		`eventKind = ?, parentConstruct = ?, version = ?, implicit = ?, deviceType = ?, deviceNumber = ?, threadId = ?, async = ?, asyncMap = ?, lineNo = ?, endLineNo = ?, funcLineNo = ?, funcEndLineNo = ?, start = ?, end = ?, cuDeviceId = ?, cuContextId = ?, cuStreamId = ?, cuProcessId = ?, cuThreadId = ?, externalId = ?, srcFile = ?, funcName = ?, numGangs = ?, numWorkers = ?, vectorLength = ?, kernelName = ?` +
		` WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, cakol.Eventkind, cakol.Parentconstruct, cakol.Version, cakol.Implicit, cakol.Devicetype, cakol.Devicenumber, cakol.Threadid, cakol.Async, cakol.Asyncmap, cakol.Lineno, cakol.Endlineno, cakol.Funclineno, cakol.Funcendlineno, cakol.Start, cakol.End, cakol.Cudeviceid, cakol.Cucontextid, cakol.Custreamid, cakol.Cuprocessid, cakol.Cuthreadid, cakol.Externalid, cakol.Srcfile, cakol.Funcname, cakol.Numgangs, cakol.Numworkers, cakol.Vectorlength, cakol.Kernelname, cakol.ID)
	_, err = db.Exec(sqlstr, cakol.Eventkind, cakol.Parentconstruct, cakol.Version, cakol.Implicit, cakol.Devicetype, cakol.Devicenumber, cakol.Threadid, cakol.Async, cakol.Asyncmap, cakol.Lineno, cakol.Endlineno, cakol.Funclineno, cakol.Funcendlineno, cakol.Start, cakol.End, cakol.Cudeviceid, cakol.Cucontextid, cakol.Custreamid, cakol.Cuprocessid, cakol.Cuthreadid, cakol.Externalid, cakol.Srcfile, cakol.Funcname, cakol.Numgangs, cakol.Numworkers, cakol.Vectorlength, cakol.Kernelname, cakol.ID)
	return err
}

// Save saves the CuptiActivityKindOpenaccLaunch to the database.
func (cakol *CuptiActivityKindOpenaccLaunch) Save(db XODB) error {
	if cakol.Exists() {
		return cakol.Update(db)
	}

	return cakol.Insert(db)
}

// Delete deletes the CuptiActivityKindOpenaccLaunch from the database.
func (cakol *CuptiActivityKindOpenaccLaunch) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !cakol._exists {
		return nil
	}

	// if deleted, bail
	if cakol._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM CUPTI_ACTIVITY_KIND_OPENACC_LAUNCH WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, cakol.ID)
	_, err = db.Exec(sqlstr, cakol.ID)
	if err != nil {
		return err
	}

	// set deleted
	cakol._deleted = true

	return nil
}

// CuptiActivityKindOpenaccLaunchByID retrieves a row from 'CUPTI_ACTIVITY_KIND_OPENACC_LAUNCH' as a CuptiActivityKindOpenaccLaunch.
//
// Generated from index 'CUPTI_ACTIVITY_KIND_OPENACC_LAUNCH__id__pkey'.
func CuptiActivityKindOpenaccLaunchByID(db XODB, id sql.NullInt64) (*CuptiActivityKindOpenaccLaunch, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`_id_, eventKind, parentConstruct, version, implicit, deviceType, deviceNumber, threadId, async, asyncMap, lineNo, endLineNo, funcLineNo, funcEndLineNo, start, end, cuDeviceId, cuContextId, cuStreamId, cuProcessId, cuThreadId, externalId, srcFile, funcName, numGangs, numWorkers, vectorLength, kernelName ` +
		`FROM CUPTI_ACTIVITY_KIND_OPENACC_LAUNCH ` +
		`WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, id)
	cakol := CuptiActivityKindOpenaccLaunch{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&cakol.ID, &cakol.Eventkind, &cakol.Parentconstruct, &cakol.Version, &cakol.Implicit, &cakol.Devicetype, &cakol.Devicenumber, &cakol.Threadid, &cakol.Async, &cakol.Asyncmap, &cakol.Lineno, &cakol.Endlineno, &cakol.Funclineno, &cakol.Funcendlineno, &cakol.Start, &cakol.End, &cakol.Cudeviceid, &cakol.Cucontextid, &cakol.Custreamid, &cakol.Cuprocessid, &cakol.Cuthreadid, &cakol.Externalid, &cakol.Srcfile, &cakol.Funcname, &cakol.Numgangs, &cakol.Numworkers, &cakol.Vectorlength, &cakol.Kernelname)
	if err != nil {
		return nil, err
	}

	return &cakol, nil
}

// CuptiActivityKindOpenaccLaunchesByStart retrieves a row from 'CUPTI_ACTIVITY_KIND_OPENACC_LAUNCH' as a CuptiActivityKindOpenaccLaunch.
//
// Generated from index 'INDEX_CUPTI_ACTIVITY_KIND_OPENACC_LAUNCH'.
func CuptiActivityKindOpenaccLaunchesByStart(db XODB, start int) ([]*CuptiActivityKindOpenaccLaunch, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`_id_, eventKind, parentConstruct, version, implicit, deviceType, deviceNumber, threadId, async, asyncMap, lineNo, endLineNo, funcLineNo, funcEndLineNo, start, end, cuDeviceId, cuContextId, cuStreamId, cuProcessId, cuThreadId, externalId, srcFile, funcName, numGangs, numWorkers, vectorLength, kernelName ` +
		`FROM CUPTI_ACTIVITY_KIND_OPENACC_LAUNCH ` +
		`WHERE start = ?`

	// run query
	XOLog(sqlstr, start)
	q, err := db.Query(sqlstr, start)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*CuptiActivityKindOpenaccLaunch{}
	for q.Next() {
		cakol := CuptiActivityKindOpenaccLaunch{
			_exists: true,
		}

		// scan
		err = q.Scan(&cakol.ID, &cakol.Eventkind, &cakol.Parentconstruct, &cakol.Version, &cakol.Implicit, &cakol.Devicetype, &cakol.Devicenumber, &cakol.Threadid, &cakol.Async, &cakol.Asyncmap, &cakol.Lineno, &cakol.Endlineno, &cakol.Funclineno, &cakol.Funcendlineno, &cakol.Start, &cakol.End, &cakol.Cudeviceid, &cakol.Cucontextid, &cakol.Custreamid, &cakol.Cuprocessid, &cakol.Cuthreadid, &cakol.Externalid, &cakol.Srcfile, &cakol.Funcname, &cakol.Numgangs, &cakol.Numworkers, &cakol.Vectorlength, &cakol.Kernelname)
		if err != nil {
			return nil, err
		}

		res = append(res, &cakol)
	}

	return res, nil
}
