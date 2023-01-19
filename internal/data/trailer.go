package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type Trailer struct {
	ID       int64    `json:"id"`
	Name     string   `json:"name"`
	Duration string   `json:"duration"`
	Genres   []string `json:"genres"`
	Date     string   `json:"date"`
}

type TrailerModel struct {
	DB *sql.DB
}

func (t TrailerModel) Insert(trailer *Trailer) error {
	query := `
	INSERT INTO trailers (name, duration, genres, date)
	VALUES ($1, $2, $3, $4)
	RETURNING id`
	args := []any{trailer.Name, trailer.Duration, pq.Array(trailer.Genres), trailer.Date}
	return t.DB.QueryRow(query, args...).Scan(&trailer.ID)
}

func (t TrailerModel) GetByGenre(genre string) ([]*Trailer, error) {
	query := fmt.Sprintf(`
	SELECT id, name, duration, genres, date
	FROM trailers
	WHERE $1 = any(genres)`)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := t.DB.QueryContext(ctx, query, genre)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	trailers := []*Trailer{}

	for rows.Next() {
		// Initialize an empty Movie struct to hold the data for an individual movie.
		var trailer Trailer
		err := rows.Scan(
			&trailer.ID,
			&trailer.Name,
			&trailer.Duration,
			pq.Array(&trailer.Genres),
			&trailer.Date,
		)
		if err != nil {
			return nil, err
		}
		trailers = append(trailers, &trailer)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return trailers, nil
}

func (t TrailerModel) GetUniv(name string, duration string, genres string) ([]*Trailer, error) {

	query := fmt.Sprintf(`
	SELECT id, name, duration, genres, date
	FROM trailers
	WHERE ((to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '') 
	AND (to_tsvector('simple', duration) @@ plainto_tsquery('simple', $2) OR $2 = ''))`)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := t.DB.QueryContext(ctx, query, name, duration)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	trailers := []*Trailer{}

	for rows.Next() {
		// Initialize an empty Movie struct to hold the data for an individual movie.
		var trailer Trailer
		err := rows.Scan(
			&trailer.ID,
			&trailer.Name,
			&trailer.Duration,
			pq.Array(&trailer.Genres),
			&trailer.Date,
		)
		if err != nil {
			return nil, err
		}
		trailers = append(trailers, &trailer)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return trailers, nil
}
