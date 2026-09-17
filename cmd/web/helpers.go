package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/nosurf"
)

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(&dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
	}
	return err
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	if app.debug {
		http.Error(w, trace, http.StatusInternalServerError)
		return
	}

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	trace := string(debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	trace := string(debug.Stack())
	app.errorLog.Output(2, trace)
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the tempalate %s does not exist", page)
		app.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.WriteHeader(status)

	buf.WriteTo(w)
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.IsAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
		Query:           app.sessionManager.PopString(r.Context(), "query"),
		Theme:           themeFromRequest(r),
	}
}

func (app *application) IsAuthenticated(r *http.Request) bool {
	IsAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return IsAuthenticated
}

type envelope map[string]any

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

func (app *application)        writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	
	js, err := json.MarshalIndent(data, " ", "\t")
	if err != nil {
		return err
	}
	js = append(js, '\n')
	
	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
			case errors.As(err, &syntaxError):
				return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
			
			
			
			case errors.Is(err, io.ErrUnexpectedEOF):
				return errors.New("body contains badly-formed JSON")

			case errors.As(err, &unmarshalTypeError):
				if unmarshalTypeError.Field != "" {
					return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
			
		
			case strings.HasPrefix(err.Error(), "json: unknown field "):
				fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
				return fmt.Errorf("body contains unknown key %s", fieldName)
				case errors.As(err, &maxBytesError):
				return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)
				case errors.As(err, &invalidUnmarshalError):
					panic(err)
				default:
					return err
				}
			}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}
