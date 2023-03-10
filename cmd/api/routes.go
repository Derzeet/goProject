package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.listMoviesHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.updateMovieHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.deleteMovieHandler)

	router.HandlerFunc(http.MethodPost, "/v1/directors", app.createDirectorHandler)
	router.HandlerFunc(http.MethodGet, "/v1/directors", app.showByAny)

	router.HandlerFunc(http.MethodPost, "/v1/trailers", app.createTrailerHandler)
	router.HandlerFunc(http.MethodGet, "/v1/trailers", app.searchTrailers)
	// router.HandlerFunc(http.MethodGet, "/v1/directors", app.listDirectorsHandler)
	// router.HandlerFunc(http.MethodPost, "/v1/actors", app.createActorHandler)
	// router.HandlerFunc(http.MethodGet, "/v1/actors/:id", app.showActorHandler)

	// Return the httprouter instance.
	return router
}
