package jellyfin

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aksiksi/gelatin/api"
)

type JellyfinApiClient struct {
	client   *http.Client
	hostname string
}

func NewJellyfinApiClient(hostname string, client *http.Client) *JellyfinApiClient {
	return &JellyfinApiClient{
		client:   client,
		hostname: hostname,
	}
}

func (c *JellyfinApiClient) GetVersion() (string, error) {
	resp, err := c.SystemInfoPublic()
	if err != nil {
		return "", err
	}

	return resp.Version, nil
}

func (c *JellyfinApiClient) SystemPing() error {
	url := fmt.Sprintf("%s%s", c.hostname, jellyfinSystemPingEndpoint)
	resp, err := c.client.Get(url)
	if err != nil {
		return err
	}

	if err := api.HttpStatusToErr(resp.StatusCode); err != nil {
		return err
	}

	return nil
}

func (c *JellyfinApiClient) SystemInfo() (*JellyfinSystemInfoResponse, error) {
	url := fmt.Sprintf("%s%s", c.hostname, jellyfinSystemInfoEndpoint)
	raw, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}

	if err := api.HttpStatusToErr(raw.StatusCode); err != nil {
		return nil, err
	}

	resp := &JellyfinSystemInfoResponse{}
	dec := json.NewDecoder(raw.Body)
	if err := dec.Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *JellyfinApiClient) SystemInfoPublic() (*JellyfinSystemInfoPublicResponse, error) {
	url := fmt.Sprintf("%s%s", c.hostname, jellyfinSystemInfoPublicEndpoint)
	raw, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}

	if err := api.HttpStatusToErr(raw.StatusCode); err != nil {
		return nil, err
	}

	resp := &JellyfinSystemInfoPublicResponse{}
	dec := json.NewDecoder(raw.Body)
	if err := dec.Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}
