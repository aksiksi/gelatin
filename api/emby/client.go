package emby

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func (c *EmbyApiClient) request(method string, url string, body io.Reader, key *api.ApiKey) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if key != nil {
		req.Header.Add(embyApiKeyHeaderName, key.ToString())
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

func (c *EmbyApiClient) get(url string, key *api.ApiKey) (*http.Response, error) {
	resp, err := c.request(http.MethodGet, url, nil, key)
	if err != nil {
		return nil, err
	}

	return resp, nil
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
	_, err := c.request(http.MethodPost, url, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *EmbyApiClient) SystemLogs(key api.ApiKey, name string) (io.ReadCloser, error) {
	url := fmt.Sprintf("%s%s/%s", c.hostname, embySystemLogsEndpoint, name)
	resp, err := c.get(url, &key)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (c *EmbyApiClient) SystemLogsQuery(key api.ApiKey) (*EmbySystemLogsQueryResponse, error) {
	url := fmt.Sprintf("%s%s", c.hostname, embySystemLogsQueryEndpoint)
	raw, err := c.get(url, &key)
	if err != nil {
		return nil, err
	}

	resp := &EmbySystemLogsQueryResponse{}
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

func (c *EmbyApiClient) SystemInfo(key api.ApiKey) (*EmbySystemInfoResponse, error) {
	url := fmt.Sprintf("%s%s", c.hostname, embySystemInfoEndpoint)
	raw, err := c.get(url, &key)
	if err != nil {
		return nil, err
	}

	resp := &EmbySystemInfoResponse{}
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

func (c *EmbyApiClient) SystemInfoPublic() (*EmbySystemInfoPublicResponse, error) {
	url := fmt.Sprintf("%s%s", c.hostname, embySystemInfoPublicEndpoint)
	raw, err := c.get(url, nil)
	if err != nil {
		return nil, err
	}

	resp := &EmbySystemInfoPublicResponse{}
	dec := json.NewDecoder(raw.Body)
	if err := dec.Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *EmbyApiClient) UserQueryPublic() (*EmbyUserQueryResponse, error) {
	url := fmt.Sprintf("%s%s", c.hostname, embyUserQueryPublicEndpoint)
	raw, err := c.get(url, nil)
	if err != nil {
		return nil, err
	}

	resp := &EmbyUserQueryResponse{}
	dec := json.NewDecoder(raw.Body)
	if err := dec.Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *EmbyApiClient) UserQuery(key api.ApiKey) (*EmbyUserQueryResponse, error) {
	url := fmt.Sprintf("%s%s", c.hostname, embyUserQueryEndpoint)
	raw, err := c.get(url, &key)
	if err != nil {
		return nil, err
	}

	resp := &EmbyUserQueryResponse{}
	dec := json.NewDecoder(raw.Body)
	if err := dec.Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *EmbyApiClient) UserGet(key api.ApiKey, userId string) (*EmbyUserDto, error) {
	url := fmt.Sprintf("%s%s/%s", c.hostname, embyUserGetEndpoint, userId)
	raw, err := c.get(url, &key)
	if err != nil {
		return nil, err
	}

	resp := &EmbyUserDto{}
	dec := json.NewDecoder(raw.Body)
	if err := dec.Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *EmbyApiClient) UserUpdate(key api.ApiKey, userId string, dto *EmbyUserDto) error {
	url := fmt.Sprintf("%s%s/%s", c.hostname, embyUserUpdateEndpoint, userId)

	data, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	_, err = c.request(http.MethodPost, url, bytes.NewReader(data), &key)
	if err != nil {
		return err
	}

	return nil
}

func (c *EmbyApiClient) UserNew(key api.ApiKey, name string) (*EmbyUserDto, error) {
	type createUserByName struct {
		Name string
	}
	req := &createUserByName{Name: name}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s%s", c.hostname, embyUserNewEndpoint)
	raw, err := c.request(http.MethodPost, url, bytes.NewReader(data), &key)
	if err != nil {
		return nil, err
	}

	resp := &EmbyUserDto{}
	dec := json.NewDecoder(raw.Body)
	if err := dec.Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}
