package nvprof

import (
	"database/sql"
	"fmt"

	"github.com/Unknwon/com"
	"github.com/pkg/errors"
	"github.com/xo/dburl"
)

type NVProf struct {
	db   *sql.DB
	path string
}

func NewNVProf(path string) (*NVProf, error) {
	if !com.IsFile(path) {
		return nil, errors.Errorf("the file at %s does not exist", path)
	}
	db, err := dburl.Open(fmt.Sprintf("file:%s?loc=auto", path))
	if err != nil {
		return nil, errors.Wrapf(err, "cannot open database at %s", path)
	}

	return &NVProf{
		db:   db,
		path: path,
	}, nil
}
