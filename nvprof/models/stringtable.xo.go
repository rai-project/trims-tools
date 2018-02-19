// Package models contains the types for schema ''.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
	"errors"
)

// Stringtable represents a row from 'StringTable'.
type Stringtable struct {
	ID    sql.NullInt64 `json:"_id_"`  // _id_
	Value string        `json:"value"` // value

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the Stringtable exists in the database.
func (s *Stringtable) Exists() bool {
	return s._exists
}

// Deleted provides information if the Stringtable has been deleted from the database.
func (s *Stringtable) Deleted() bool {
	return s._deleted
}

// Insert inserts the Stringtable to the database.
func (s *Stringtable) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if s._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO StringTable (` +
		`value` +
		`) VALUES (` +
		`?` +
		`)`

	// run query
	XOLog(sqlstr, s.Value)
	res, err := db.Exec(sqlstr, s.Value)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	s.ID = sql.NullInt64(id)
	s._exists = true

	return nil
}

// Update updates the Stringtable in the database.
func (s *Stringtable) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !s._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if s._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE StringTable SET ` +
		`value = ?` +
		` WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, s.Value, s.ID)
	_, err = db.Exec(sqlstr, s.Value, s.ID)
	return err
}

// Save saves the Stringtable to the database.
func (s *Stringtable) Save(db XODB) error {
	if s.Exists() {
		return s.Update(db)
	}

	return s.Insert(db)
}

// Delete deletes the Stringtable from the database.
func (s *Stringtable) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !s._exists {
		return nil
	}

	// if deleted, bail
	if s._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM StringTable WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, s.ID)
	_, err = db.Exec(sqlstr, s.ID)
	if err != nil {
		return err
	}

	// set deleted
	s._deleted = true

	return nil
}

// StringtableByID retrieves a row from 'StringTable' as a Stringtable.
//
// Generated from index 'StringTable__id__pkey'.
func StringtableByID(db XODB, id sql.NullInt64) (*Stringtable, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`_id_, value ` +
		`FROM StringTable ` +
		`WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, id)
	s := Stringtable{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&s.ID, &s.Value)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// StringtableByValue retrieves a row from 'StringTable' as a Stringtable.
//
// Generated from index 'sqlite_autoindex_StringTable_1'.
func StringtableByValue(db XODB, value string) (*Stringtable, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`_id_, value ` +
		`FROM StringTable ` +
		`WHERE value = ?`

	// run query
	XOLog(sqlstr, value)
	s := Stringtable{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, value).Scan(&s.ID, &s.Value)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
