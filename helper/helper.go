package helper

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

type APIHandler func(http.ResponseWriter, *http.Request) error

func MakeHTTPHandler(h APIHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			var res APIError
			if errors.As(err, &res) {
				WriteJSON(w, res.StatusCode, res)
			} else {
				res = APIError{
					StatusCode: http.StatusInternalServerError,
					Msg:        "Internal Server Error",
				}
				WriteJSON(w, res.StatusCode, res)
			}
			slog.Error("REST API", "err", err.Error(), "path", r.URL.Path)
		}
	}
}

func ParseJSON(r *http.Request, payload any) error {
	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")

	if tokenAuth != "" {
		return tokenAuth
	}

	if tokenQuery != "" {
		return tokenQuery
	}

	return ""
}
