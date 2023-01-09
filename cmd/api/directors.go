package main

import (
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

// func (app *application) searchByNameHandler(w http.ResponseWriter, r *http.Request) {
// 	qs := r.URL.Query()
// 	name := app.readString(qs, "name", "")

// 	director, err := app.models.Directors.(name)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, data.ErrRecordNotFound):
// 			app.notFoundResponse(w, r)
// 		default:
// 			app.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}
// 	err = app.writeJSON(w, http.StatusOK, envelope{"director": director}, nil)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}

// }

func (app *application) listDirectorsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string
		Surname  string
		Awards   []string
		FullName string
		data.Filters
	}
	qs := r.URL.Query()
	input.Name = app.readString(qs, "name", "")
	input.Surname = app.readString(qs, "surname", "")
	input.Awards = app.readCSV(qs, "awards", []string{})
	input.FullName = app.readString(qs, "fullname", "")
	input.Filters.Page = app.readInt(qs, "page", 1)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "surname", "-id", "-name", "-awards", "awards", "-surname", "-runtime"}

	// Call the GetAll() method to retrieve the movies, passing in the various filter // parameters.
	directors, err := app.models.Directors.GetAllDirectors(input.Name, input.Surname, input.Awards, input.Filters, input.FullName)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Send a JSON response containing the movie data.
	err = app.writeJSON(w, http.StatusOK, envelope{"director": directors}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
