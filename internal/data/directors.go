package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
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

func (d DirectorModel) Insert(director *Director) error {
	query := `
	INSERT INTO directors (name, surname, awards)
	VALUES ($1, $2, $3)
	RETURNING id`
	args := []any{director.Name, director.Surname, pq.Array(director.Awards)}
	return d.DB.QueryRow(query, args...).Scan(&director.ID)
}

func (d DirectorModel) GetByName(name string) (*Director, error) {
	var director Director
	if name == "" {
		return nil, ErrRecordNotFound
	}
	query := `SELECT id, name, surname, awards FROM directors WHERE name = $1`

	result := d.DB.QueryRow(query, name).Scan(
		&director.ID,
		&director.Name,
		&director.Surname,
		pq.Array(&director.Awards),
	)

	if result != nil {
		switch {
		case errors.Is(result, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, result
		}
	}

	return &director, nil
}

// func (d DirectorModel) getByAw(award string) []*Director {
// 	directors := []*Director{}
// 	query := `select * from directors where $1 = any(awards)`

// 	row, err := d.DB.Query(query, award)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer row.Close()
// 	for row.Next() {
// 		var director Director
// 		err := row.Scan(
// 			&director.ID,
// 			&director.Name,
// 			&director.Surname,
// 			pq.Array(&director.Awards),
// 		)
// 		if err != nil {
// 			panic(err)
// 		}
// 		directors = append(directors, &director)
// 	}

// 	return directors
// }

func (d DirectorModel) GetUniv(name string, surname string, award string) ([]*Director, error) {

	query := fmt.Sprintf(`
	SELECT id, name, surname, awards
	FROM directors
	WHERE ((to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '') 
	AND (to_tsvector('simple', surname) @@ plainto_tsquery('simple', $2) OR $2 = '')) 
	OR $3 = any(awards)`)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := d.DB.QueryContext(ctx, query, name, surname, award)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	directors := []*Director{}

	for rows.Next() {
		// Initialize an empty Movie struct to hold the data for an individual movie.
		var director Director
		err := rows.Scan(
			&director.ID,
			&director.Name,
			&director.Surname,
			pq.Array(&director.Awards),
		)
		if err != nil {
			return nil, err
		}
		directors = append(directors, &director)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return directors, nil
}
