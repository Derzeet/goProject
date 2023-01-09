package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Movies MovieModel
	// Actors ActorModel
	Directors DirectorModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies:    MovieModel{DB: db},
		Directors: DirectorModel{DB: db},
	}
}

func (d DirectorModel) Insert(director *Director) error {
	query := `
	INSERT INTO directors (name, surname, awards)
	VALUES ($1, $2, $3)
	RETURNING id`
	args := []any{director.Name, director.Surname, pq.Array(director.Awards)}
	return d.DB.QueryRow(query, args...).Scan(&director.ID)
}

// func (d DirectorModel) searchBy(name string) (*Director, error) {
// 	query := `select id, name, surname, awards from directors where name = $1`
// 	var director Director
// 	err := d.DB.QueryRow(query, name).Scan(
// 		&director.ID,
// 		&director.Name,
// 		&director.Surname,
// 		pq.Array(&director.Awards),
// 	)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, sql.ErrNoRows):
// 			return nil, ErrRecordNotFound
// 		default:
// 			return nil, err
// 		}
// 	}
// 	return &director, nil
// }

func (d DirectorModel) GetAllDirectors(name string, surname string, awards []string, filters Filters, fullname string) ([]*Director, error) { // Update the SQL query to include the filter conditions.
	query := fmt.Sprintf(`
SELECT id, name, surname, awards
FROM directors
WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '') AND (awards @> $2 OR $2 = '{}') AND (to_tsvector('simple', surname) @@ plainto_tsquery('simple', $5) OR $5 = ''))
ORDER BY %s %s, id ASC
LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []any{name, pq.Array(awards), filters.limit(), filters.offset(), surname}
	rows, err := d.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	director := []*Director{}
	for rows.Next() {
		var directors Director
		err := rows.Scan(&directors.ID,
			&directors.Name, &directors.Surname, pq.Array(&directors.Awards),
		)
		if err != nil {
			return nil, err
		}
		director = append(director, &directors)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return director, nil
}
