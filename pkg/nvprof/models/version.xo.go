// Package models contains the types for schema ''.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
)

// Version represents a row from 'Version'.
type Version struct {
	Version sql.NullInt64 `json:"version"` // version
}