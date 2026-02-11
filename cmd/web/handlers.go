package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	

	"github.com/julienschmidt/httprouter"
	"snippetbox.net/internal/models"
	"snippetbox.net/internal/validator"
)

type snippetCreateForm struct {
    Title       string `form:"title"`
    Content     string `form:"content"`
    Expires     int	   `form:"expires"`
	validator.Validator`form:"-"`
}




func (app *application) home(w http.ResponseWriter, r *http.Request) {

	res, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = res


	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
    params := httprouter.ParamsFromContext(r.Context())

    id, err := strconv.Atoi(params.ByName("id"))
    if err != nil || id < 1 {
        app.notFound(w)
        return
    }

    snippet, err := app.snippets.Get(id)
    if err != nil {
        if errors.Is(err, models.ErrNoRecord) {
            app.notFound(w)
        } else {
            app.serverError(w, err)
        }
        return
    }

    data := app.newTemplateData(r)
    data.Snippet = snippet

    app.render(w, http.StatusOK, "view.tmpl", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {


    data := app.newTemplateData(r)

	data.Form = snippetCreateForm{
        Expires: 365,
    }

    app.render(w, http.StatusOK, "create.tmpl", data)


}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {



	var form snippetCreateForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}



	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChar(form.Title, 100), "title", "This field cannot be more 100 characters long")

	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.MaxChar(form.Content, 1000), "content", "This field cannot be more 1000 characters long")

	form.CheckField(validator.Accept_values(form.Expires, 1, 7, 365), "expires", "Expires cannot be current value")


    if !form.Valid() {
        data := app.newTemplateData(r)
        data.Form = form
        app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
        return
    }

    res, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
    if err != nil {
        app.serverError(w, err)
        return
    }

    app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", res), http.StatusSeeOther)
}
