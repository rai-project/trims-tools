// Package models contains the types for schema ''.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
	"errors"
)

// CuptiActivityKindEventInstance represents a row from 'CUPTI_ACTIVITY_KIND_EVENT_INSTANCE'.
type CuptiActivityKindEventInstance struct {
	ID            sql.NullInt64 `json:"_id_"`          // _id_
	ID            int           `json:"id"`            // id
	Domain        int           `json:"domain"`        // domain
	Instance      int           `json:"instance"`      // instance
	Value         int           `json:"value"`         // value
	Correlationid int           `json:"correlationId"` // correlationId

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the CuptiActivityKindEventInstance exists in the database.
func (cakei *CuptiActivityKindEventInstance) Exists() bool {
	return cakei._exists
}

// Deleted provides information if the CuptiActivityKindEventInstance has been deleted from the database.
func (cakei *CuptiActivityKindEventInstance) Deleted() bool {
	return cakei._deleted
}

// Insert inserts the CuptiActivityKindEventInstance to the database.
func (cakei *CuptiActivityKindEventInstance) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if cakei._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO CUPTI_ACTIVITY_KIND_EVENT_INSTANCE (` +
		`domain, instance, value, correlationId` +
		`) VALUES (` +
		`?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, cakei.Domain, cakei.Instance, cakei.Value, cakei.Correlationid)
	res, err := db.Exec(sqlstr, cakei.Domain, cakei.Instance, cakei.Value, cakei.Correlationid)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	cakei.ID = sql.NullInt64(id)
	cakei._exists = true

	return nil
}

// Update updates the CuptiActivityKindEventInstance in the database.
func (cakei *CuptiActivityKindEventInstance) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !cakei._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if cakei._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE CUPTI_ACTIVITY_KIND_EVENT_INSTANCE SET ` +
		`domain = ?, instance = ?, value = ?, correlationId = ?` +
		` WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, cakei.Domain, cakei.Instance, cakei.Value, cakei.Correlationid, cakei.ID)
	_, err = db.Exec(sqlstr, cakei.Domain, cakei.Instance, cakei.Value, cakei.Correlationid, cakei.ID)
	return err
}

// Save saves the CuptiActivityKindEventInstance to the database.
func (cakei *CuptiActivityKindEventInstance) Save(db XODB) error {
	if cakei.Exists() {
		return cakei.Update(db)
	}

	return cakei.Insert(db)
}

// Delete deletes the CuptiActivityKindEventInstance from the database.
func (cakei *CuptiActivityKindEventInstance) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !cakei._exists {
		return nil
	}

	// if deleted, bail
	if cakei._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM CUPTI_ACTIVITY_KIND_EVENT_INSTANCE WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, cakei.ID)
	_, err = db.Exec(sqlstr, cakei.ID)
	if err != nil {
		return err
	}

	// set deleted
	cakei._deleted = true

	return nil
}

// CuptiActivityKindEventInstanceByID retrieves a row from 'CUPTI_ACTIVITY_KIND_EVENT_INSTANCE' as a CuptiActivityKindEventInstance.
//
// Generated from index 'CUPTI_ACTIVITY_KIND_EVENT_INSTANCE__id__pkey'.
func CuptiActivityKindEventInstanceByID(db XODB, id sql.NullInt64) (*CuptiActivityKindEventInstance, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`_id_, id, domain, instance, value, correlationId ` +
		`FROM CUPTI_ACTIVITY_KIND_EVENT_INSTANCE ` +
		`WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, id)
	cakei := CuptiActivityKindEventInstance{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&cakei.ID, &cakei.ID, &cakei.Domain, &cakei.Instance, &cakei.Value, &cakei.Correlationid)
	if err != nil {
		return nil, err
	}

	return &cakei, nil
}

// CuptiActivityKindEventInstancesByCorrelationid retrieves a row from 'CUPTI_ACTIVITY_KIND_EVENT_INSTANCE' as a CuptiActivityKindEventInstance.
//
// Generated from index 'INDEX_CUPTI_ACTIVITY_KIND_EVENT_INSTANCE'.
func CuptiActivityKindEventInstancesByCorrelationid(db XODB, correlationid int) ([]*CuptiActivityKindEventInstance, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`_id_, id, domain, instance, value, correlationId ` +
		`FROM CUPTI_ACTIVITY_KIND_EVENT_INSTANCE ` +
		`WHERE correlationId = ?`

	// run query
	XOLog(sqlstr, correlationid)
	q, err := db.Query(sqlstr, correlationid)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*CuptiActivityKindEventInstance{}
	for q.Next() {
		cakei := CuptiActivityKindEventInstance{
			_exists: true,
		}

		// scan
		err = q.Scan(&cakei.ID, &cakei.ID, &cakei.Domain, &cakei.Instance, &cakei.Value, &cakei.Correlationid)
		if err != nil {
			return nil, err
		}

		res = append(res, &cakei)
	}

	return res, nil
}
