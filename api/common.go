package api

import (
	"errors"
	"fmt"
	"net/http"
)

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
