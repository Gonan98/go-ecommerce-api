package helper

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type APIError struct {
	StatusCode int `json:"statusCode"`
	Msg        any `json:"error"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("API error: %d", e.StatusCode)
}

func InvalidJSON() APIError {
	return APIError{
		StatusCode: http.StatusBadRequest,
		Msg:        "Invalid JSON request data",
	}
}

func BadRequest(msg string) APIError {
	return APIError{
		StatusCode: http.StatusBadRequest,
		Msg:        msg,
	}
}

func PermissionDenied() APIError {
	return APIError{
		StatusCode: http.StatusForbidden,
		Msg:        "Permission denied",
	}
}

func InvalidEmailOrPassword() APIError {
	return APIError{
		StatusCode: http.StatusBadRequest,
		Msg:        "Invalid email of password",
	}
}

func InvalidRequestData(errors validator.ValidationErrors) APIError {

	mapped := make(map[string]string)

	for _, e := range errors {
		s := strings.Split(e.Error(), ":")
		mapped[e.Field()] = s[2]
	}

	return APIError{
		StatusCode: http.StatusBadRequest,
		Msg:        mapped,
	}
}

func EntityNotFound(entity string, id int) APIError {
	return APIError{
		StatusCode: http.StatusNotFound,
		Msg:        fmt.Sprintf("%s with ID=%d does not exist", entity, id),
	}
}
