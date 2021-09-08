package jellyfin

import (
	"encoding/json"
	"fmt"
	"io"
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

func (c *JellyfinApiClient) request(method string, url string, body io.Reader, key *api.ApiKey) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if key != nil {
		req.Header.Add(jellyfinApiKeyHeaderName, key.ToString())
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if err := api.HttpStatusToErr(resp.StatusCode); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *JellyfinApiClient) get(url string, key *api.ApiKey) (*http.Response, error) {
	resp, err := c.request(http.MethodGet, url, nil, key)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *JellyfinApiClient) SystemPing() error {
	url := fmt.Sprintf("%s%s", c.hostname, jellyfinSystemPingEndpoint)
	resp, err := c.get(url, nil)
	if err != nil {
		return err
	}

	if err := api.HttpStatusToErr(resp.StatusCode); err != nil {
		return err
	}

	return nil
}

func (c *JellyfinApiClient) SystemLogs(key api.ApiKey, name string) (io.ReadCloser, error) {
	url := fmt.Sprintf("%s%s/%s", c.hostname, jellyfinSystemLogsEndpoint, name)
	resp, err := c.get(url, &key)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (c *JellyfinApiClient) SystemLogsQuery(key api.ApiKey) (*JellyfinSystemLogsQueryResponse, error) {
	url := fmt.Sprintf("%s%s", c.hostname, jellyfinSystemLogsQueryEndpoint)
	raw, err := c.get(url, &key)
	if err != nil {
		return nil, err
	}

	resp := &JellyfinSystemLogsQueryResponse{}
	dec := json.NewDecoder(raw.Body)
	if err := dec.Decode(resp); err != nil {
		return nil, err
	}

	err = api.Validator.Struct(resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *JellyfinApiClient) SystemInfo(key api.ApiKey) (*JellyfinSystemInfoResponse, error) {
	url := fmt.Sprintf("%s%s", c.hostname, jellyfinSystemInfoEndpoint)
	raw, err := c.get(url, &key)
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

	err = api.Validator.Struct(resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *JellyfinApiClient) SystemInfoPublic() (*JellyfinSystemInfoPublicResponse, error) {
	url := fmt.Sprintf("%s%s", c.hostname, jellyfinSystemInfoPublicEndpoint)
	raw, err := c.get(url, nil)
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
