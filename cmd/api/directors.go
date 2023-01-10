package main

import (
	"errors"
	"fmt"
	"net/http"

	"assOneGo.derzeet.net/internal/data"
)

func (app *application) createDirectorHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name    string   `json:"name"`
		Surname string   `json:"surname"`
		Awards  []string `json:"awards"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	director := &data.Director{
		Name:    input.Name,
		Surname: input.Surname,
		Awards:  input.Awards,
	}
	err = app.models.Directors.Insert(director)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/directors/%d", director.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"director": director}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) showDirector(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	name := app.readString(qs, "name", "")

	director, err := app.models.Directors.GetByName(name)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"director": director}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showByAny(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	award := app.readString(qs, "award", "")
	name := app.readString(qs, "name", "")
	surname := app.readString(qs, "surname", "")
	directors, err := app.models.Directors.GetUniv(name, surname, award)
	if err != nil {
		panic(err)
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"directors": directors}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
