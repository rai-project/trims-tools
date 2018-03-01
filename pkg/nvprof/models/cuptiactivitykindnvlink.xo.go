// Package models contains the types for schema ''.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
	"errors"
)

// CuptiActivityKindNvlink represents a row from 'CUPTI_ACTIVITY_KIND_NVLINK'.
type CuptiActivityKindNvlink struct {
	ID                  sql.NullInt64 `json:"_id_"`                // _id_
	Nvlinkversion       int           `json:"nvlinkVersion"`       // nvlinkVersion
	Typedev0            int           `json:"typeDev0"`            // typeDev0
	Typedev1            int           `json:"typeDev1"`            // typeDev1
	Iddev0              []byte        `json:"idDev0"`              // idDev0
	Iddev1              []byte        `json:"idDev1"`              // idDev1
	Flag                int           `json:"flag"`                // flag
	Physicalnvlinkcount int           `json:"physicalNvLinkCount"` // physicalNvLinkCount
	Portdev0            []byte        `json:"portDev0"`            // portDev0
	Portdev1            []byte        `json:"portDev1"`            // portDev1
	Bandwidth           int           `json:"bandwidth"`           // bandwidth

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the CuptiActivityKindNvlink exists in the database.
func (cakn *CuptiActivityKindNvlink) Exists() bool {
	return cakn._exists
}

// Deleted provides information if the CuptiActivityKindNvlink has been deleted from the database.
func (cakn *CuptiActivityKindNvlink) Deleted() bool {
	return cakn._deleted
}

// Insert inserts the CuptiActivityKindNvlink to the database.
func (cakn *CuptiActivityKindNvlink) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if cakn._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO CUPTI_ACTIVITY_KIND_NVLINK (` +
		`nvlinkVersion, typeDev0, typeDev1, idDev0, idDev1, flag, physicalNvLinkCount, portDev0, portDev1, bandwidth` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, cakn.Nvlinkversion, cakn.Typedev0, cakn.Typedev1, cakn.Iddev0, cakn.Iddev1, cakn.Flag, cakn.Physicalnvlinkcount, cakn.Portdev0, cakn.Portdev1, cakn.Bandwidth)
	res, err := db.Exec(sqlstr, cakn.Nvlinkversion, cakn.Typedev0, cakn.Typedev1, cakn.Iddev0, cakn.Iddev1, cakn.Flag, cakn.Physicalnvlinkcount, cakn.Portdev0, cakn.Portdev1, cakn.Bandwidth)
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

// Update updates the CuptiActivityKindNvlink in the database.
func (cakn *CuptiActivityKindNvlink) Update(db XODB) error {
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
	const sqlstr = `UPDATE CUPTI_ACTIVITY_KIND_NVLINK SET ` +
		`nvlinkVersion = ?, typeDev0 = ?, typeDev1 = ?, idDev0 = ?, idDev1 = ?, flag = ?, physicalNvLinkCount = ?, portDev0 = ?, portDev1 = ?, bandwidth = ?` +
		` WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, cakn.Nvlinkversion, cakn.Typedev0, cakn.Typedev1, cakn.Iddev0, cakn.Iddev1, cakn.Flag, cakn.Physicalnvlinkcount, cakn.Portdev0, cakn.Portdev1, cakn.Bandwidth, cakn.ID)
	_, err = db.Exec(sqlstr, cakn.Nvlinkversion, cakn.Typedev0, cakn.Typedev1, cakn.Iddev0, cakn.Iddev1, cakn.Flag, cakn.Physicalnvlinkcount, cakn.Portdev0, cakn.Portdev1, cakn.Bandwidth, cakn.ID)
	return err
}

// Save saves the CuptiActivityKindNvlink to the database.
func (cakn *CuptiActivityKindNvlink) Save(db XODB) error {
	if cakn.Exists() {
		return cakn.Update(db)
	}

	return cakn.Insert(db)
}

// Delete deletes the CuptiActivityKindNvlink from the database.
func (cakn *CuptiActivityKindNvlink) Delete(db XODB) error {
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
	const sqlstr = `DELETE FROM CUPTI_ACTIVITY_KIND_NVLINK WHERE _id_ = ?`

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

// CuptiActivityKindNvlinkByID retrieves a row from 'CUPTI_ACTIVITY_KIND_NVLINK' as a CuptiActivityKindNvlink.
//
// Generated from index 'CUPTI_ACTIVITY_KIND_NVLINK__id__pkey'.
func CuptiActivityKindNvlinkByID(db XODB, id sql.NullInt64) (*CuptiActivityKindNvlink, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`_id_, nvlinkVersion, typeDev0, typeDev1, idDev0, idDev1, flag, physicalNvLinkCount, portDev0, portDev1, bandwidth ` +
		`FROM CUPTI_ACTIVITY_KIND_NVLINK ` +
		`WHERE _id_ = ?`

	// run query
	XOLog(sqlstr, id)
	cakn := CuptiActivityKindNvlink{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&cakn.ID, &cakn.Nvlinkversion, &cakn.Typedev0, &cakn.Typedev1, &cakn.Iddev0, &cakn.Iddev1, &cakn.Flag, &cakn.Physicalnvlinkcount, &cakn.Portdev0, &cakn.Portdev1, &cakn.Bandwidth)
	if err != nil {
		return nil, err
	}

	return &cakn, nil
}