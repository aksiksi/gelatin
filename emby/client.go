package emby

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	gelatin "github.com/aksiksi/gelatin/lib"
)

const (
	embyApiKeyAuthHeader  = "X-Emby-Authorization"
	embyApiKeyTokenHeader = "X-Emby-Token"
)

const (
	embySystemPingEndpoint       = "/System/Ping"
	embySystemLogsEndpoint       = "/System/Logs"
	embySystemLogsQueryEndpoint  = "/System/Logs/Query"
	embySystemInfoEndpoint       = "/System/Info"
	embySystemInfoPublicEndpoint = "/System/Info/Public"
	embyUserQueryEndpoint        = "/Users/Query"
	embyUserQueryPublicEndpoint  = "/Users/Public"
	embyUserGetEndpoint          = "/Users"
	embyUserUpdateEndpoint       = "/Users"
	embyUserNewEndpoint          = "/Users/New"
	embyUserDeleteEndpoint       = "/Users"
	embyUserPasswordEndpoint     = "/Users"
	embyUserAuthEndpoint         = "/Users/AuthenticateByName"
	embyUserPolicyEndpoint       = "/Users"
)

type embyApiKey struct {
	key string
}

// NewApiKey returns a new ApiKey for the given client
func NewApiKey(key string) gelatin.ApiKey {
	return &embyApiKey{
		key: key,
	}
}

func (k *embyApiKey) ToString() string {
	return k.key
}

func (*embyApiKey) HeaderName() string {
	return embyApiKeyTokenHeader
}

type EmbyApiClient struct {
	client   *http.Client
	hostname string
}

func NewEmbyApiClient(hostname string, client *http.Client) *EmbyApiClient {
	return &EmbyApiClient{
		client:   client,
		hostname: fmt.Sprintf("%s/emby", hostname),
	}
}

func (c *EmbyApiClient) System() gelatin.GelatinSystemService {
	// TODO: Move this out
	return c
}

func (c *EmbyApiClient) User() gelatin.GelatinUserService {
	// TODO: Move this out
	return c
}

func (c *EmbyApiClient) request(method string, url string, body io.Reader, key gelatin.ApiKey) (*http.Response, error) {
	headers := map[string]string{
		embyApiKeyAuthHeader: `Emby Client="gelatin", Device="gelatin", DeviceId="007", Version="0.0.1"`,
	}

	if key != nil {
		headers[embyApiKeyTokenHeader] = key.ToString()
	}

	if body != nil {
		headers["Content-Type"] = "application/json"
	}

	resp, err := gelatin.HttpRequest(c.client, method, url, body, headers)

	return resp, err
}

func (c *EmbyApiClient) get(url string, key gelatin.ApiKey) (*http.Response, error) {
	resp, err := c.request(http.MethodGet, url, nil, key)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *EmbyApiClient) Version() (string, error) {
	info, err := c.Info(nil, true)
	if err != nil {
		return "", err
	}

	return info.Version, nil
}

func (c *EmbyApiClient) Ping() error {
	url := fmt.Sprintf("%s%s", c.hostname, embySystemPingEndpoint)
	_, err := c.request(http.MethodPost, url, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *EmbyApiClient) GetLogs(key gelatin.ApiKey) ([]gelatin.GelatinSystemLog, error) {
	url := fmt.Sprintf("%s%s", c.hostname, embySystemLogsQueryEndpoint)
	raw, err := c.get(url, key)
	if err != nil {
		return nil, err
	}

	resp := &EmbySystemLogsQueryResponse{}
	dec := json.NewDecoder(raw.Body)
	if err := dec.Decode(resp); err != nil {
		return nil, err
	}

	if resp.TotalRecordCount <= 0 {
		return nil, fmt.Errorf("invalid record count")
	}

	return resp.Items, nil
}

func (c *EmbyApiClient) GetLogFile(key gelatin.ApiKey, name string) (io.ReadCloser, error) {
	url := fmt.Sprintf("%s%s/%s", c.hostname, embySystemLogsEndpoint, name)
	resp, err := c.get(url, key)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (c *EmbyApiClient) Info(key gelatin.ApiKey, public bool) (*gelatin.GelatinSystemInfo, error) {
	var url string

	if public {
		url = fmt.Sprintf("%s%s", c.hostname, embySystemInfoPublicEndpoint)
	} else {
		url = fmt.Sprintf("%s%s", c.hostname, embySystemInfoEndpoint)
	}

	raw, err := c.get(url, key)
	if err != nil {
		return nil, err
	}

	resp := &gelatin.GelatinSystemInfo{}
	dec := json.NewDecoder(raw.Body)
	if err := dec.Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *EmbyApiClient) GetUser(key gelatin.ApiKey, id string) (*gelatin.GelatinUser, error) {
	url := fmt.Sprintf("%s%s/%s", c.hostname, embyUserGetEndpoint, id)
	raw, err := c.get(url, key)
	if err != nil {
		return nil, err
	}

	resp := &gelatin.GelatinUser{}
	dec := json.NewDecoder(raw.Body)
	if err := dec.Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *EmbyApiClient) GetUsers(key gelatin.ApiKey, public bool) ([]gelatin.GelatinUser, error) {
	var url string

	if public {
		url = fmt.Sprintf("%s%s", c.hostname, embyUserQueryPublicEndpoint)
	} else {
		url = fmt.Sprintf("%s%s", c.hostname, embyUserQueryEndpoint)
	}

	raw, err := c.get(url, key)
	if err != nil {
		return nil, err
	}

	if public {
		var resp []gelatin.GelatinUser

		dec := json.NewDecoder(raw.Body)
		if err := dec.Decode(&resp); err != nil {
			return nil, err
		}

		return resp, nil
	} else {
		resp := &EmbyUserQueryResponse{}
		dec := json.NewDecoder(raw.Body)
		if err := dec.Decode(&resp); err != nil {
			return nil, err
		}

		if resp.TotalRecordCount <= 0 {
			return nil, fmt.Errorf("invalid record count")
		}

		return resp.Items, nil
	}
}

func (c *EmbyApiClient) UpdateUser(key gelatin.AdminKey, id string, data *gelatin.GelatinUser) error {
	url := fmt.Sprintf("%s%s/%s", c.hostname, embyUserUpdateEndpoint, id)

	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = c.request(http.MethodPost, url, bytes.NewReader(raw), key)
	if err != nil {
		return err
	}

	return nil
}

func (c *EmbyApiClient) NewUser(key gelatin.AdminKey, name string) (*gelatin.GelatinUser, error) {
	type createUserByName struct {
		Name string
	}
	req := &createUserByName{Name: name}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s%s", c.hostname, embyUserNewEndpoint)
	raw, err := c.request(http.MethodPost, url, bytes.NewReader(data), key)
	if err != nil {
		return nil, err
	}

	resp := &gelatin.GelatinUser{}
	dec := json.NewDecoder(raw.Body)
	if err := dec.Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *EmbyApiClient) DeleteUser(key gelatin.AdminKey, id string) error {
	url := fmt.Sprintf("%s%s/%s", c.hostname, embyUserDeleteEndpoint, id)

	_, err := c.request(http.MethodDelete, url, nil, key)
	if err != nil {
		return err
	}

	return nil
}

func (c *EmbyApiClient) UpdatePassword(key gelatin.AdminKey, id, currentPassword, newPassword string, reset bool) error {
	type setUserPassword struct {
		Id        string
		CurrentPw string
		NewPw     string
		Reset     bool
	}

	req := setUserPassword{
		Id:        id,
		CurrentPw: currentPassword,
		NewPw:     newPassword,
		Reset:     reset,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s%s/%s/Password", c.hostname, embyUserPasswordEndpoint, id)
	_, err = c.request(http.MethodPost, url, bytes.NewReader(data), key)
	if err != nil {
		return err
	}

	return nil
}

func (c *EmbyApiClient) Authenticate(username, password string) (userKey gelatin.ApiKey, err error) {
	req := map[string]string{
		"Username": username,
		"Pw":       password,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s%s", c.hostname, embyUserAuthEndpoint)
	raw, err := c.request(http.MethodPost, url, bytes.NewReader(data), nil)
	if err != nil {
		return nil, err
	}

	resp := &EmbyUserAuthResponse{}

	dec := json.NewDecoder(raw.Body)
	if err := dec.Decode(resp); err != nil {
		return nil, err
	}

	return NewApiKey(resp.AccessToken), nil
}

func (c *EmbyApiClient) UpdatePolicy(key gelatin.AdminKey, id string, policy *gelatin.GelatinUserPolicy) error {
	url := fmt.Sprintf("%s%s/%s/Policy", c.hostname, embyUserPolicyEndpoint, id)

	data, err := json.Marshal(policy)
	if err != nil {
		return err
	}

	_, err = c.request(http.MethodPost, url, bytes.NewReader(data), key)
	if err != nil {
		return err
	}

	return nil
}
