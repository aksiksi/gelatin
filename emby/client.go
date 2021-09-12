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

func (c *EmbyApiClient) request(method string, url string, body io.Reader, key gelatin.ApiKey) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add(embyApiKeyAuthHeader, `Emby Client="gelatin", Device="gelatin", DeviceId="007", Version="0.0.1"`)

	if key != nil {
		req.Header.Add(embyApiKeyTokenHeader, key.ToString())
	}

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if err := gelatin.HttpStatusToErr(resp.StatusCode); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *EmbyApiClient) get(url string, key gelatin.ApiKey) (*http.Response, error) {
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

func (c *EmbyApiClient) SystemLogs(key gelatin.ApiKey, name string) (io.ReadCloser, error) {
	url := fmt.Sprintf("%s%s/%s", c.hostname, embySystemLogsEndpoint, name)
	resp, err := c.get(url, key)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (c *EmbyApiClient) SystemLogsQuery(key gelatin.ApiKey) (*EmbySystemLogsQueryResponse, error) {
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

	return resp, nil
}

func (c *EmbyApiClient) SystemInfo(key gelatin.ApiKey) (*EmbySystemInfoResponse, error) {
	url := fmt.Sprintf("%s%s", c.hostname, embySystemInfoEndpoint)
	raw, err := c.get(url, key)
	if err != nil {
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

func (c *EmbyApiClient) UserQueryPublic() ([]*EmbyUserDto, error) {
	url := fmt.Sprintf("%s%s", c.hostname, embyUserQueryPublicEndpoint)
	raw, err := c.get(url, nil)
	if err != nil {
		return nil, err
	}

	var resp []*EmbyUserDto
	dec := json.NewDecoder(raw.Body)
	if err := dec.Decode(&resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *EmbyApiClient) UserQuery(key gelatin.ApiKey) (*EmbyUserQueryResponse, error) {
	url := fmt.Sprintf("%s%s", c.hostname, embyUserQueryEndpoint)
	raw, err := c.get(url, key)
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

func (c *EmbyApiClient) UserGet(key gelatin.ApiKey, userId string) (*EmbyUserDto, error) {
	url := fmt.Sprintf("%s%s/%s", c.hostname, embyUserGetEndpoint, userId)
	raw, err := c.get(url, key)
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

func (c *EmbyApiClient) UserUpdate(key gelatin.ApiKey, userId string, dto *EmbyUserDto) error {
	url := fmt.Sprintf("%s%s/%s", c.hostname, embyUserUpdateEndpoint, userId)

	data, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	_, err = c.request(http.MethodPost, url, bytes.NewReader(data), key)
	if err != nil {
		return err
	}

	return nil
}

func (c *EmbyApiClient) UserNew(key gelatin.ApiKey, name string) (*EmbyUserDto, error) {
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

	resp := &EmbyUserDto{}
	dec := json.NewDecoder(raw.Body)
	if err := dec.Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *EmbyApiClient) UserDelete(key gelatin.ApiKey, userId string) error {
	url := fmt.Sprintf("%s%s/%s", c.hostname, embyUserDeleteEndpoint, userId)

	_, err := c.request(http.MethodDelete, url, nil, key)
	if err != nil {
		return err
	}

	return nil
}

func (c *EmbyApiClient) ResetUserPassword(key gelatin.ApiKey, userId string) error {
	type resetUserPassword struct {
		Id            string
		ResetPassword bool
	}

	req := resetUserPassword{Id: userId, ResetPassword: true}
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s%s/%s/Password", c.hostname, embyUserPasswordEndpoint, userId)
	_, err = c.request(http.MethodPost, url, bytes.NewReader(data), key)
	if err != nil {
		return err
	}

	return nil
}

func (c *EmbyApiClient) UserPassword(key gelatin.ApiKey, userId, currentPassword, newPassword string, reset bool) error {
	type setUserPassword struct {
		Id        string
		CurrentPw string
		NewPw     string
		Reset     bool
	}

	req := setUserPassword{
		Id:        userId,
		CurrentPw: currentPassword,
		NewPw:     newPassword,
		Reset:     reset,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s%s/%s/Password", c.hostname, embyUserPasswordEndpoint, userId)
	_, err = c.request(http.MethodPost, url, bytes.NewReader(data), key)
	if err != nil {
		return err
	}

	return nil
}

func (c *EmbyApiClient) UserAuth(username, password string) (userKey gelatin.ApiKey, err error) {
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

func (c *EmbyApiClient) UserPolicy(key gelatin.AdminKey, userId string, policy *EmbyUserPolicy) error {
	url := fmt.Sprintf("%s%s/%s/Policy", c.hostname, embyUserPolicyEndpoint, userId)

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
