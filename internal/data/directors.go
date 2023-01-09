package data

import (
	"database/sql"
)

type Director struct {
	ID      int64    `json:"id"`
	Name    string   `json:"name"`
	Surname string   `json:"surname"`
	Awards  []string `json:"awards"`
}

type DirectorModel struct {
	DB *sql.DB
}
