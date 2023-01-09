package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

type Actor struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Age       int32     `json:"age"`
	Movies    []int32   `json:"movies"`
	Version   int32     `json:"version"`
}

type ActorModel struct {
	DB *sql.DB
}

func (a ActorModel) Insert(actor *Actor) error {
	query := `
	INSERT INTO actors (first_name, last_name, age, movies)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, version`
	args := []any{actor.FirstName, actor.LastName, actor.Age, pq.Array(actor.Movies)}
	return a.DB.QueryRow(query, args...).Scan(&actor.ID, &actor.CreatedAt, &actor.Version)
}

func (a ActorModel) Get(id int64) (*Actor, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `SELECT id, created_at, first_name, last_name, age, movies, version from actors where id = $1`
	var actor Actor
	err := a.DB.QueryRow(query, id).Scan(
		&actor.ID,
		&actor.CreatedAt,
		&actor.FirstName,
		&actor.LastName,
		&actor.Age,
		pq.Array(&actor.Movies),
		&actor.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &actor, nil
}
