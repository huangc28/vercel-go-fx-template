package render

import (
	"errors"
	"net/http"

	chirender "github.com/go-chi/render"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func ChiJSON(w http.ResponseWriter, r *http.Request, status int, v any) {
	chirender.Status(r, status)
	chirender.JSON(w, r, v)
}

func ChiErr(w http.ResponseWriter, r *http.Request, status int, err error) {
	if err == nil {
		err = errors.New("unknown error")
	}
	ChiJSON(w, r, status, ErrorResponse{Error: err.Error()})
}
