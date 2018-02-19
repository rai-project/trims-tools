// Package models contains the types for schema ''.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
	"errors"
)

// CuptiActivityKindName represents a row from 'CUPTI_ACTIVITY_KIND_NAME'.
type CuptiActivityKindName struct {
	ID         sql.NullInt64 `json:"_id_"`       // _id_
	Objectkind int           `json:"objectKind"` // objectKind
	Objectid   []byte        `json:"objectId"`   // objectId
	Name       int           `json:"name"`       // name

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the CuptiActivityKindName exists in the database.
func (cakn *CuptiActivityKindName) Exists() bool {
	return cakn._exists
}

// Deleted provides information if the CuptiActivityKindName has been deleted from the database.
func (cakn *CuptiActivityKindName) Deleted() bool {
	return cakn._deleted
}

// Insert inserts the CuptiActivityKindName to the database.
func (cakn *CuptiActivityKindName) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if cakn._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO CUPTI_ACTIVITY_KIND_NAME (` +
		`objectKind, objectId, name` +
		`) VALUES (` +
		`?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, cakn.Objectkind, cakn.Objectid, cakn.Name)
	res, err := db.Exec(sqlstr, cakn.Objectkind, cakn.Objectid, cakn.Name)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	cakn.ID = sql.NullInt64(id)
	cakn._exists = true

	return nil
}

// Update updates the CuptiActivityKindName in the database.
func (cakn *CuptiActivityKindName) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !cakn._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if cakn._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE CUPTI_ACTIVITY_KIND_NAME SET ` +
		`objectKind = ?, objectId = ?, name = ?` +
		` WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, cakn.Objectkind, cakn.Objectid, cakn.Name, cakn.ID)
	_, err = db.Exec(sqlstr, cakn.Objectkind, cakn.Objectid, cakn.Name, cakn.ID)
	return err
}

// Save saves the CuptiActivityKindName to the database.
func (cakn *CuptiActivityKindName) Save(db XODB) error {
	if cakn.Exists() {
		return cakn.Update(db)
	}

	return cakn.Insert(db)
}

// Delete deletes the CuptiActivityKindName from the database.
func (cakn *CuptiActivityKindName) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !cakn._exists {
		return nil
	}

	// if deleted, bail
	if cakn._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM CUPTI_ACTIVITY_KIND_NAME WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, cakn.ID)
	_, err = db.Exec(sqlstr, cakn.ID)
	if err != nil {
		return err
	}

	// set deleted
	cakn._deleted = true

	return nil
}

// CuptiActivityKindNameByID retrieves a row from 'CUPTI_ACTIVITY_KIND_NAME' as a CuptiActivityKindName.
//
// Generated from index 'CUPTI_ACTIVITY_KIND_NAME__id__pkey'.
func CuptiActivityKindNameByID(db XODB, id sql.NullInt64) (*CuptiActivityKindName, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`_id_, objectKind, objectId, name ` +
		`FROM CUPTI_ACTIVITY_KIND_NAME ` +
		`WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, id)
	cakn := CuptiActivityKindName{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&cakn.ID, &cakn.Objectkind, &cakn.Objectid, &cakn.Name)
	if err != nil {
		return nil, err
	}

	return &cakn, nil
}
