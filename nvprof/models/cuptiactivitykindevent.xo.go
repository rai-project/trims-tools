// Package models contains the types for schema ''.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
	"errors"
)

// CuptiActivityKindEvent represents a row from 'CUPTI_ACTIVITY_KIND_EVENT'.
type CuptiActivityKindEvent struct {
	ID            sql.NullInt64 `json:"_id_"`          // _id_
	ID            int           `json:"id"`            // id
	Value         int           `json:"value"`         // value
	Domain        int           `json:"domain"`        // domain
	Correlationid int           `json:"correlationId"` // correlationId

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the CuptiActivityKindEvent exists in the database.
func (cake *CuptiActivityKindEvent) Exists() bool {
	return cake._exists
}

// Deleted provides information if the CuptiActivityKindEvent has been deleted from the database.
func (cake *CuptiActivityKindEvent) Deleted() bool {
	return cake._deleted
}

// Insert inserts the CuptiActivityKindEvent to the database.
func (cake *CuptiActivityKindEvent) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if cake._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO CUPTI_ACTIVITY_KIND_EVENT (` +
		`value, domain, correlationId` +
		`) VALUES (` +
		`?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, cake.Value, cake.Domain, cake.Correlationid)
	res, err := db.Exec(sqlstr, cake.Value, cake.Domain, cake.Correlationid)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	cake.ID = sql.NullInt64(id)
	cake._exists = true

	return nil
}

// Update updates the CuptiActivityKindEvent in the database.
func (cake *CuptiActivityKindEvent) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !cake._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if cake._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE CUPTI_ACTIVITY_KIND_EVENT SET ` +
		`value = ?, domain = ?, correlationId = ?` +
		` WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, cake.Value, cake.Domain, cake.Correlationid, cake.ID)
	_, err = db.Exec(sqlstr, cake.Value, cake.Domain, cake.Correlationid, cake.ID)
	return err
}

// Save saves the CuptiActivityKindEvent to the database.
func (cake *CuptiActivityKindEvent) Save(db XODB) error {
	if cake.Exists() {
		return cake.Update(db)
	}

	return cake.Insert(db)
}

// Delete deletes the CuptiActivityKindEvent from the database.
func (cake *CuptiActivityKindEvent) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !cake._exists {
		return nil
	}

	// if deleted, bail
	if cake._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM CUPTI_ACTIVITY_KIND_EVENT WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, cake.ID)
	_, err = db.Exec(sqlstr, cake.ID)
	if err != nil {
		return err
	}

	// set deleted
	cake._deleted = true

	return nil
}

// CuptiActivityKindEventByID retrieves a row from 'CUPTI_ACTIVITY_KIND_EVENT' as a CuptiActivityKindEvent.
//
// Generated from index 'CUPTI_ACTIVITY_KIND_EVENT__id__pkey'.
func CuptiActivityKindEventByID(db XODB, id sql.NullInt64) (*CuptiActivityKindEvent, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`_id_, id, value, domain, correlationId ` +
		`FROM CUPTI_ACTIVITY_KIND_EVENT ` +
		`WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, id)
	cake := CuptiActivityKindEvent{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&cake.ID, &cake.ID, &cake.Value, &cake.Domain, &cake.Correlationid)
	if err != nil {
		return nil, err
	}

	return &cake, nil
}

// CuptiActivityKindEventsByCorrelationid retrieves a row from 'CUPTI_ACTIVITY_KIND_EVENT' as a CuptiActivityKindEvent.
//
// Generated from index 'INDEX_CUPTI_ACTIVITY_KIND_EVENT'.
func CuptiActivityKindEventsByCorrelationid(db XODB, correlationid int) ([]*CuptiActivityKindEvent, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`_id_, id, value, domain, correlationId ` +
		`FROM CUPTI_ACTIVITY_KIND_EVENT ` +
		`WHERE correlationId = ?`

	// run query
	XOLog(sqlstr, correlationid)
	q, err := db.Query(sqlstr, correlationid)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*CuptiActivityKindEvent{}
	for q.Next() {
		cake := CuptiActivityKindEvent{
			_exists: true,
		}

		// scan
		err = q.Scan(&cake.ID, &cake.ID, &cake.Value, &cake.Domain, &cake.Correlationid)
		if err != nil {
			return nil, err
		}

		res = append(res, &cake)
	}

	return res, nil
}
