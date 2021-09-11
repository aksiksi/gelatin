package api

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validator *validator.Validate = validator.New()

type ApiKey interface {
	// String representation of the key for use in a header
	ToString() string
	// Name of the header to use
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
