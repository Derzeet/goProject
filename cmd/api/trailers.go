package main

import (
	"fmt"
	"net/http"

	"assOneGo.derzeet.net/internal/data"
)

func (app *application) createTrailerHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string   `json:"name"`
		Duration string   `json:"duration"`
		Genres   []string `json:"genres"`
		Date     string   `json:"date"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	trailer := &data.Trailer{
		Name:     input.Name,
		Duration: input.Duration,
		Genres:   input.Genres,
		Date:     input.Date,
	}
	err = app.models.Trailers.Insert(trailer)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/directors/%d", trailer.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"director": trailer}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) searchTrailers(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	genre := app.readString(qs, "genre", "")
	if genre != "" {
		trailers, err := app.models.Trailers.GetByGenre(genre)
		if err != nil {
			panic(err)
		}
		err = app.writeJSON(w, http.StatusOK, envelope{"trailers": trailers}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	} else {
		name := app.readString(qs, "name", "")
		duration := app.readString(qs, "date", "")
		trailers, err := app.models.Trailers.GetUniv(name, duration, genre)
		if err != nil {
			panic(err)
		}
		err = app.writeJSON(w, http.StatusOK, envelope{"trailers": trailers}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}
}
