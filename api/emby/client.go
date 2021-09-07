package emby

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aksiksi/gelatin/api"
)

type EmbyApiClient struct {
	client   *http.Client
	hostname string
}

func NewEmbyApiClient(hostname string, client *http.Client) *EmbyApiClient {
	return &EmbyApiClient{
		client:   client,
		hostname: hostname,
	}
}

func (c *EmbyApiClient) GetVersion() (string, error) {
	resp, err := c.SystemInfoPublic()
	if err != nil {
		return "", err
	}

	return resp.Version, nil
}

func (c *EmbyApiClient) SystemPing() error {
	url := fmt.Sprintf("%s%s", c.hostname, embySystemPingEndpoint)
	resp, err := c.client.Post(url, "text/plain", nil)
	if err != nil {
		return err
	}

	if err := api.HttpStatusToErr(resp.StatusCode); err != nil {
		return err
	}

	return nil
}

func (c *EmbyApiClient) SystemInfo() (*EmbySystemInfoResponse, error) {
	url := fmt.Sprintf("%s%s", c.hostname, embySystemInfoEndpoint)
	raw, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}

	if err := api.HttpStatusToErr(raw.StatusCode); err != nil {
		return nil, err
	}

	resp := &EmbySystemInfoResponse{}
	dec := json.NewDecoder(raw.Body)
	if err := dec.Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *EmbyApiClient) SystemInfoPublic() (*EmbySystemInfoPublicResponse, error) {
	url := fmt.Sprintf("%s%s", c.hostname, embySystemInfoPublicEndpoint)
	raw, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}

	if err := api.HttpStatusToErr(raw.StatusCode); err != nil {
		return nil, err
	}

	resp := &EmbySystemInfoPublicResponse{}
	dec := json.NewDecoder(raw.Body)
	if err := dec.Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}
