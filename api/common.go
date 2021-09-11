package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validator *validator.Validate = validator.New()

type ApiKey struct {
	key      string
	client   string
	device   string
	deviceId string
	version  string
}

// NewApiKey returns a new ApiKey
//
// If auth is true, the key will contain full authorization info
func NewApiKey(key string, auth bool) ApiKey {
	if !auth {
		return ApiKey{key: key}
	} else {
		return ApiKey{
			key:      key,
			client:   "gelatin",
			device:   "gelatin",
			deviceId: "gelatin",
			version:  "0.0.1",
		}
	}
}

func (k ApiKey) ToString() string {
	if k.client == "" {
		return k.key
	} else {
		return fmt.Sprintf(
			`MediaBrowser Client="%s", Device="%s", DeviceId="%s", Version="%s", Token="%s"`,
			k.client, k.device, k.deviceId, k.version, k.key,
		)
	}
}

func HttpStatusToErr(code int) error {
	switch code {
	case http.StatusOK, http.StatusNoContent:
		return nil
	default:
		return errors.New(http.StatusText(code))
	}
}
