// Package models contains the types for schema ''.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
	"errors"
)

// CuptiActivityKindSourceLocator represents a row from 'CUPTI_ACTIVITY_KIND_SOURCE_LOCATOR'.
type CuptiActivityKindSourceLocator struct {
	ID         sql.NullInt64 `json:"_id_"`       // _id_
	ID         int           `json:"id"`         // id
	Linenumber int           `json:"lineNumber"` // lineNumber
	Filename   int           `json:"fileName"`   // fileName

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the CuptiActivityKindSourceLocator exists in the database.
func (caksl *CuptiActivityKindSourceLocator) Exists() bool {
	return caksl._exists
}

// Deleted provides information if the CuptiActivityKindSourceLocator has been deleted from the database.
func (caksl *CuptiActivityKindSourceLocator) Deleted() bool {
	return caksl._deleted
}

// Insert inserts the CuptiActivityKindSourceLocator to the database.
func (caksl *CuptiActivityKindSourceLocator) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if caksl._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO CUPTI_ACTIVITY_KIND_SOURCE_LOCATOR (` +
		`lineNumber, fileName` +
		`) VALUES (` +
		`?, ?` +
		`)`

	// run query
	XOLog(sqlstr, caksl.Linenumber, caksl.Filename)
	res, err := db.Exec(sqlstr, caksl.Linenumber, caksl.Filename)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	caksl.ID = sql.NullInt64(id)
	caksl._exists = true

	return nil
}

// Update updates the CuptiActivityKindSourceLocator in the database.
func (caksl *CuptiActivityKindSourceLocator) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !caksl._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if caksl._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE CUPTI_ACTIVITY_KIND_SOURCE_LOCATOR SET ` +
		`lineNumber = ?, fileName = ?` +
		` WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, caksl.Linenumber, caksl.Filename, caksl.ID)
	_, err = db.Exec(sqlstr, caksl.Linenumber, caksl.Filename, caksl.ID)
	return err
}

// Save saves the CuptiActivityKindSourceLocator to the database.
func (caksl *CuptiActivityKindSourceLocator) Save(db XODB) error {
	if caksl.Exists() {
		return caksl.Update(db)
	}

	return caksl.Insert(db)
}

// Delete deletes the CuptiActivityKindSourceLocator from the database.
func (caksl *CuptiActivityKindSourceLocator) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !caksl._exists {
		return nil
	}

	// if deleted, bail
	if caksl._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM CUPTI_ACTIVITY_KIND_SOURCE_LOCATOR WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, caksl.ID)
	_, err = db.Exec(sqlstr, caksl.ID)
	if err != nil {
		return err
	}

	// set deleted
	caksl._deleted = true

	return nil
}

// CuptiActivityKindSourceLocatorByID retrieves a row from 'CUPTI_ACTIVITY_KIND_SOURCE_LOCATOR' as a CuptiActivityKindSourceLocator.
//
// Generated from index 'CUPTI_ACTIVITY_KIND_SOURCE_LOCATOR__id__pkey'.
func CuptiActivityKindSourceLocatorByID(db XODB, id sql.NullInt64) (*CuptiActivityKindSourceLocator, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`_id_, id, lineNumber, fileName ` +
		`FROM CUPTI_ACTIVITY_KIND_SOURCE_LOCATOR ` +
		`WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, id)
	caksl := CuptiActivityKindSourceLocator{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&caksl.ID, &caksl.ID, &caksl.Linenumber, &caksl.Filename)
	if err != nil {
		return nil, err
	}

	return &caksl, nil
}