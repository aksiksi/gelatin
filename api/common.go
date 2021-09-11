package api

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validator *validator.Validate = validator.New()

type ApiKey interface {
	ToString() string
	HeaderName() string
}

func HttpStatusToErr(code int) error {
	switch code {
	case http.StatusOK, http.StatusNoContent:
		return nil
	default:
		return errors.New(http.StatusText(code))
	}
}
