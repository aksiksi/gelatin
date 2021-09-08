package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validator *validator.Validate = validator.New()

type ApiKey struct {
	key string
}

func NewApiKey(key string) ApiKey {
	return ApiKey{key}
}

func (k ApiKey) ToString() string {
	return k.key
}

func HttpStatusToErr(code int) error {
	switch code {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest:
		return errors.New("bad request")
	case http.StatusUnauthorized:
		return errors.New("unauthorized")
	case http.StatusForbidden:
		return errors.New("forbidden")
	case http.StatusNotFound:
		return errors.New("not found")
	default:
		return fmt.Errorf("internal server error: %d", code)
	}
}
