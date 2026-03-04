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

type usersSugnipForm struct {
    Name string `form:"name"`
    Email string `form:"email"`
    Password string `form:"password"`
    validator.Validator`form:"-"`
}

type usersLoginForm struct {
    Email string `form:"email"`
    Password string `form:"password"`
    validator.Validator`form:"-"`
}

type usersPasswordForm struct {
    Password string `form:"curr_password"`
    NewPassword string `form:"new_password"`
    ConfirmNewPassword string `form:"confirm_password"`
    validator.Validator`form:"-"`
}

type idSnippetForm struct {
    ID int `form:"id"`
}

func (app *application) deleteSnippetPost(w http.ResponseWriter, r *http.Request) {
    var form idSnippetForm
    err := app.decodePostForm(r, &form)
    if err != nil {
        app.clientError(w, http.StatusBadRequest)
        return
    }
    app.snippets.Delete(form.ID)
    app.sessionManager.Put(r.Context(), "flash", "Snippet delete successful")
    http.Redirect(w, r, "/", http.StatusSeeOther)
}


func (app *application) updateSnippet(w http.ResponseWriter, r *http.Request) {
   // w.Write([]byte("test"))
    params := httprouter.ParamsFromContext(r.Context())

    id, err := strconv.Atoi(params.ByName("id"))
    if err != nil || id < 1 {
        app.notFound(w)
        return
    }

    userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
    app.sessionManager.Put(r.Context(), "snippetID", id)

    snippet, err := app.snippets.Get(id, userID)
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

	data.Form = snippetCreateForm{
        Expires: 365,
        Title: snippet.Title,
        Content: snippet.Content,
    }

    app.render(w, http.StatusOK, "updateSnippet.tmpl", data)

}

func (app *application) updateSnippetPost(w http.ResponseWriter, r *http.Request) {

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
    snippetID := app.sessionManager.PopInt(r.Context(), "snippetID")
    err = app.snippets.Update(form.Title, form.Content, form.Expires, snippetID)
    if err != nil {
        app.serverError(w, err)
        return
    }
    app.sessionManager.Put(r.Context(), "flash", "Snippet successfully updated!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", snippetID), http.StatusSeeOther)
}


func (app *application) passwordUpdate(w http.ResponseWriter, r *http.Request) {
    data := app.newTemplateData(r)
    data.Form = usersPasswordForm{}
    app.render(w, http.StatusOK, "updatePassword.tmpl", data)
}

func (app *application) passwordUpdatePost(w http.ResponseWriter, r *http.Request) {
    var form usersPasswordForm
    err := app.decodePostForm(r, &form)
    if err != nil {
        app.clientError(w, http.StatusBadRequest)
        return
    }

    user_id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	form.CheckField(validator.NotBlank(form.Password), "curr_password", "This field cannot be blank")
    form.CheckField(validator.NotBlank(form.NewPassword), "new_password", "This field cannot be blank")
    form.CheckField(validator.NotBlank(form.ConfirmNewPassword), "confirm_password", "This field cannot be blank")

    b, err := app.users.ReturnCorrectPassword(form.Password, user_id)
    if err != nil {
        if errors.Is(err, models.ErrInvalidCredentials) {
            form.CheckField(b, "curr_password", "Password is incorrect")
        } else {
            app.serverError(w, err)
        } 
    }
    
    form.CheckField(validator.MaxChar(form.NewPassword, 100), "new_password", "This field cannot be more 100 characters long")
    form.CheckField(validator.MinChar(form.NewPassword, 7), "new_password", "This field cannot be less 8 characters long")
    form.CheckField(validator.SamePassword(form.NewPassword, form.ConfirmNewPassword), "confirm_password", "the passwords don't match")

    if !form.Valid() {
        data := app.newTemplateData(r)
        data.Form = form
        app.render(w, http.StatusUnprocessableEntity, "updatePassword.tmpl", data)
        return
    }

    err = app.users.ChangePassword(form.NewPassword, user_id)
    if err != nil {
        app.serverError(w, err)
    }
    data := app.newTemplateData(r)
    data.Form = form
    app.sessionManager.Put(r.Context(), "flash", "Change password complete success")
    http.Redirect(w, r, "/account/view", http.StatusSeeOther)
}


func (app *application) account(w http.ResponseWriter, r *http.Request) {

    data := app.newTemplateData(r)
    id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
    var err error
    data.Form, err = app.users.ReturnData(id)
    if err != nil {
        if errors.Is(err, models.ErrNoRecord) {
            app.notFound(w)
        } else {
            app.serverError(w, err)
        }
        return
    }
    app.render(w, http.StatusOK, "account.tmpl", data)

}


func (app *application) about(w http.ResponseWriter, r *http.Request) {
    data := app.newTemplateData(r)
    app.render(w, http.StatusOK, "about.tmpl", data)
    
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {

    data := app.newTemplateData(r)

    data.Form = usersSugnipForm{}

    app.render(w, http.StatusOK, "signup.tmpl", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
    var form usersSugnipForm
    err := app.decodePostForm(r, &form)
    if err != nil {
        app.clientError(w, http.StatusBadRequest)
        return
    }

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
    form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
    form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")

	form.CheckField(validator.MaxChar(form.Name, 32), "name", "This field cannot be more 33 characters long")
	form.CheckField(validator.MaxChar(form.Password, 100), "password", "This field cannot be more 100 characters long")
	form.CheckField(validator.MaxChar(form.Email, 100), "email", "This field cannot be more 100 characters long")

	form.CheckField(validator.MinChar(form.Password, 7), "password", "This field cannot be less 8 characters long")

    form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")

    if !form.Valid() {
        data := app.newTemplateData(r)
        data.Form = form
        app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
        return
    }

    err = app.users.Insert(form.Name, form.Email, form.Password)
    if err != nil {
        if errors.Is(err, models.ErrDuplicateEmail){

            form.AddFieldMap("email", "Email address is already in use")
            data := app.newTemplateData(r)
            data.Form = form
            app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
        } else {
            app.serverError(w, err)
        }
        return
    }
    app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")
    http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {

    data := app.newTemplateData(r)
    data.Form = usersLoginForm{}
    app.render(w, http.StatusOK, "login.tmpl", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
    var form usersLoginForm

    err := app.decodePostForm(r, &form)
    if err != nil {
        app.clientError(w, http.StatusBadRequest)
        return
    }

    form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
    form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
    form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

    if !form.Valid() {
        data := app.newTemplateData(r)
        data.Form = form
        app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
        return
    }

    id, err := app.users.Authenticate(form.Email, form.Password)
    if err != nil {
        if errors.Is(err, models.ErrInvalidCredentials) {
            form.AddNonFieldError("Email or password is incorrect")

            data := app.newTemplateData(r)
            data.Form = form
            app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
        } else {
            app.serverError(w, err)
        }
        return
    }

    err = app.sessionManager.RenewToken(r.Context())
    if err != nil {
        app.serverError(w, err)
        return
    }

    app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

    if path :=app.sessionManager.GetString(r.Context(),"RedirectPathAfterLogin"); path != "" {
    http.Redirect(w, r, path, http.StatusSeeOther)

    } else {
    http.Redirect(w, r, "/account/view", http.StatusSeeOther)
    }

}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {

    if  !app.IsAuthenticated(r) {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    err := app.sessionManager.RenewToken(r.Context())
    if err != nil {
        app.serverError(w, err)
        
    }

    app.sessionManager.Remove(r.Context(), "authenticatedUserID")
    app.sessionManager.Remove(r.Context(), "RedirectPathAfterLogin")

    app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

    http.Redirect(w, r, "/", http.StatusSeeOther)
}


func (app *application) home(w http.ResponseWriter, r *http.Request) {

    user_id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	res, err := app.snippets.Latest(user_id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = res


	app.render(w, http.StatusOK, "home.tmpl", data)
}



func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
    params := httprouter.ParamsFromContext(r.Context())

    id, err := strconv.Atoi(params.ByName("id"))
    if err != nil || id < 1 {
        app.notFound(w)
        return
    }
    user_id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
    snippet, err := app.snippets.Get(id, user_id)
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
 user_id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
    res, err := app.snippets.Insert(form.Title, form.Content, form.Expires, user_id)
    if err != nil {
        app.serverError(w, err)
        return
    }

    app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", res), http.StatusSeeOther)
}


