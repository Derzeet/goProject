package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Movies    MovieModel
	Directors DirectorModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies:    MovieModel{DB: db},
		Directors: DirectorModel{DB: db},
	}
}
