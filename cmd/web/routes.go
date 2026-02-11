package main

import (
	"net/http"
	"github.com/justinas/alice"
	"github.com/julienschmidt/httprouter"

)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	fileServer := http.FileServer(http.Dir("/home/aartchik/project/golang/snippetbox/ui/static"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        app.notFound(w)
    })

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
    router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost,"/snippet/create", dynamic.ThenFunc(app.snippetCreatePost))
	router.HandlerFunc(http.MethodGet, "/ping", app.ping)

	standart := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	
	return standart.Then(router)
}
