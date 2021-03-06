package gelatin

import (
	"errors"
	"io"
	"net/http"
)

type ApiKey interface {
	ToString() string
	IsAdmin() bool
}

func httpStatusToErr(code int) error {
	switch code {
	case http.StatusOK, http.StatusNoContent:
		return nil
	default:
		return errors.New(http.StatusText(code))
	}
}

func HttpRequest(client *http.Client, method string, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if err := httpStatusToErr(resp.StatusCode); err != nil {
		return nil, err
	}

	return resp, nil
}
