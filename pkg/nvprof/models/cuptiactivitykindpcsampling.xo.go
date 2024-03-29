// Package models contains the types for schema ''.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
	"errors"
)

// CuptiActivityKindPcSampling represents a row from 'CUPTI_ACTIVITY_KIND_PC_SAMPLING'.
type CuptiActivityKindPcSampling struct {
	ID              sql.NullInt64 `json:"_id_"`            // _id_
	Flags           int           `json:"flags"`           // flags
	Sourcelocatorid int           `json:"sourceLocatorId"` // sourceLocatorId
	Correlationid   int           `json:"correlationId"`   // correlationId
	Functionid      int           `json:"functionId"`      // functionId
	Latencysamples  int           `json:"latencySamples"`  // latencySamples
	Samples         int           `json:"samples"`         // samples
	Stallreason     int           `json:"stallReason"`     // stallReason
	Pcoffset        int           `json:"pcOffset"`        // pcOffset

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the CuptiActivityKindPcSampling exists in the database.
func (cakps *CuptiActivityKindPcSampling) Exists() bool {
	return cakps._exists
}

// Deleted provides information if the CuptiActivityKindPcSampling has been deleted from the database.
func (cakps *CuptiActivityKindPcSampling) Deleted() bool {
	return cakps._deleted
}

// Insert inserts the CuptiActivityKindPcSampling to the database.
func (cakps *CuptiActivityKindPcSampling) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if cakps._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO CUPTI_ACTIVITY_KIND_PC_SAMPLING (` +
		`flags, sourceLocatorId, correlationId, functionId, latencySamples, samples, stallReason, pcOffset` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, cakps.Flags, cakps.Sourcelocatorid, cakps.Correlationid, cakps.Functionid, cakps.Latencysamples, cakps.Samples, cakps.Stallreason, cakps.Pcoffset)
	res, err := db.Exec(sqlstr, cakps.Flags, cakps.Sourcelocatorid, cakps.Correlationid, cakps.Functionid, cakps.Latencysamples, cakps.Samples, cakps.Stallreason, cakps.Pcoffset)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	cakps.ID = sql.NullInt64(id)
	cakps._exists = true

	return nil
}

// Update updates the CuptiActivityKindPcSampling in the database.
func (cakps *CuptiActivityKindPcSampling) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !cakps._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if cakps._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE CUPTI_ACTIVITY_KIND_PC_SAMPLING SET ` +
		`flags = ?, sourceLocatorId = ?, correlationId = ?, functionId = ?, latencySamples = ?, samples = ?, stallReason = ?, pcOffset = ?` +
		` WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, cakps.Flags, cakps.Sourcelocatorid, cakps.Correlationid, cakps.Functionid, cakps.Latencysamples, cakps.Samples, cakps.Stallreason, cakps.Pcoffset, cakps.ID)
	_, err = db.Exec(sqlstr, cakps.Flags, cakps.Sourcelocatorid, cakps.Correlationid, cakps.Functionid, cakps.Latencysamples, cakps.Samples, cakps.Stallreason, cakps.Pcoffset, cakps.ID)
	return err
}

// Save saves the CuptiActivityKindPcSampling to the database.
func (cakps *CuptiActivityKindPcSampling) Save(db XODB) error {
	if cakps.Exists() {
		return cakps.Update(db)
	}

	return cakps.Insert(db)
}

// Delete deletes the CuptiActivityKindPcSampling from the database.
func (cakps *CuptiActivityKindPcSampling) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !cakps._exists {
		return nil
	}

	// if deleted, bail
	if cakps._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM CUPTI_ACTIVITY_KIND_PC_SAMPLING WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, cakps.ID)
	_, err = db.Exec(sqlstr, cakps.ID)
	if err != nil {
		return err
	}

	// set deleted
	cakps._deleted = true

	return nil
}

// CuptiActivityKindPcSamplingByID retrieves a row from 'CUPTI_ACTIVITY_KIND_PC_SAMPLING' as a CuptiActivityKindPcSampling.
//
// Generated from index 'CUPTI_ACTIVITY_KIND_PC_SAMPLING__id__pkey'.
func CuptiActivityKindPcSamplingByID(db XODB, id sql.NullInt64) (*CuptiActivityKindPcSampling, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`_id_, flags, sourceLocatorId, correlationId, functionId, latencySamples, samples, stallReason, pcOffset ` +
		`FROM CUPTI_ACTIVITY_KIND_PC_SAMPLING ` +
		`WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, id)
	cakps := CuptiActivityKindPcSampling{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&cakps.ID, &cakps.Flags, &cakps.Sourcelocatorid, &cakps.Correlationid, &cakps.Functionid, &cakps.Latencysamples, &cakps.Samples, &cakps.Stallreason, &cakps.Pcoffset)
	if err != nil {
		return nil, err
	}

	return &cakps, nil
}

// CuptiActivityKindPcSamplingsByCorrelationid retrieves a row from 'CUPTI_ACTIVITY_KIND_PC_SAMPLING' as a CuptiActivityKindPcSampling.
//
// Generated from index 'INDEX_CUPTI_ACTIVITY_KIND_PC_SAMPLING'.
func CuptiActivityKindPcSamplingsByCorrelationid(db XODB, correlationid int) ([]*CuptiActivityKindPcSampling, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`_id_, flags, sourceLocatorId, correlationId, functionId, latencySamples, samples, stallReason, pcOffset ` +
		`FROM CUPTI_ACTIVITY_KIND_PC_SAMPLING ` +
		`WHERE correlationId = ?`

	// run query
	XOLog(sqlstr, correlationid)
	q, err := db.Query(sqlstr, correlationid)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*CuptiActivityKindPcSampling{}
	for q.Next() {
		cakps := CuptiActivityKindPcSampling{
			_exists: true,
		}

		// scan
		err = q.Scan(&cakps.ID, &cakps.Flags, &cakps.Sourcelocatorid, &cakps.Correlationid, &cakps.Functionid, &cakps.Latencysamples, &cakps.Samples, &cakps.Stallreason, &cakps.Pcoffset)
		if err != nil {
			return nil, err
		}

		res = append(res, &cakps)
	}

	return res, nil
}
