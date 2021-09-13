package jellyfin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	gelatin "github.com/aksiksi/gelatin/lib"
)

const (
	// See: https://github.com/jellyfin/jellyfin/blob/8c463b9b8196530e6d417b40ed81825422cada32/Jellyfin.Server.Implementations/Security/AuthorizationContext.cs#L80
	jellyfinApiKeyAuthHeader  = "X-Emby-Authorization"
	jellyfinApiKeyTokenHeader = "X-Emby-Token"
)

const (
	jellyfinSystemPingEndpoint       = "/System/Ping"
	jellyfinSystemLogsEndpoint       = "/System/Logs"
	jellyfinSystemLogsNameEndpoint   = "/System/Logs/Log"
	jellyfinSystemInfoEndpoint       = "/System/Info"
	jellyfinSystemInfoPublicEndpoint = "/System/Info/Public"
	jellyfinUserQueryEndpoint        = "/users"
	jellyfinUserQueryPublicEndpoint  = "/users/public"
	jellyfinUserGetEndpoint          = "/users"
	jellyfinUserUpdateEndpoint       = "/users"
	jellyfinUserNewEndpoint          = "/users/new"
	jellyfinUserDeleteEndpoint       = "/users"
	jellyfinUserPasswordEndpoint     = "/Users"
	jellyfinUserAuthEndpoint         = "/Users/AuthenticateByName"
	jellyfinUserPolicyEndpoint       = "/Users"
)

type jellyfinApiKey struct {
	key string
}

// NewApiKey returns a new ApiKey for the given client
func NewApiKey(key string) gelatin.ApiKey {
	return &jellyfinApiKey{
		key: key,
	}
}

func (k *jellyfinApiKey) ToString() string {
	return k.key
}

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

func (c *JellyfinApiClient) System() gelatin.GelatinSystemService {
	// TODO: Move this out
	return c
}

func (c *JellyfinApiClient) User() gelatin.GelatinUserService {
	// TODO: Move this out
	return c
}

func (c *JellyfinApiClient) request(method string, url string, body io.Reader, key gelatin.ApiKey) (*http.Response, error) {
	headers := map[string]string{
		jellyfinApiKeyAuthHeader: `MediaBrowser Client="gelatin", Device="gelatin", DeviceId="007", Version="0.0.1"`,
	}

	if key != nil {
		headers[jellyfinApiKeyTokenHeader] = key.ToString()
	}

	if body != nil {
		headers["Content-Type"] = "application/json"
	}

	resp, err := gelatin.HttpRequest(c.client, method, url, body, headers)

	return resp, err
}

func (c *JellyfinApiClient) get(url string, key gelatin.ApiKey) (*http.Response, error) {
	resp, err := c.request(http.MethodGet, url, nil, key)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *JellyfinApiClient) Version() (string, error) {
	resp, err := c.Info(nil, true)
	if err != nil {
		return "", err
	}

	return resp.Version, nil
}

func (c *JellyfinApiClient) Ping() error {
	url := fmt.Sprintf("%s%s", c.hostname, jellyfinSystemPingEndpoint)
	_, err := c.get(url, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *JellyfinApiClient) GetLogs(key gelatin.ApiKey) ([]gelatin.GelatinSystemLog, error) {
	url := fmt.Sprintf("%s%s", c.hostname, jellyfinSystemLogsEndpoint)
	raw, err := c.get(url, key)
	if err != nil {
		return nil, err
	}

	var resp []gelatin.GelatinSystemLog
	dec := json.NewDecoder(raw.Body)
	if err := dec.Decode(&resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *JellyfinApiClient) GetLogFile(key gelatin.ApiKey, name string) (io.ReadCloser, error) {
	url := fmt.Sprintf("%s%s?name=%s", c.hostname, jellyfinSystemLogsNameEndpoint, name)

	resp, err := c.get(url, key)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (c *JellyfinApiClient) Info(key gelatin.ApiKey, public bool) (*gelatin.GelatinSystemInfo, error) {
	var url string
	if public {
		url = fmt.Sprintf("%s%s", c.hostname, jellyfinSystemInfoPublicEndpoint)
	} else {
		url = fmt.Sprintf("%s%s", c.hostname, jellyfinSystemInfoEndpoint)
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

func (c *JellyfinApiClient) GetUser(key gelatin.ApiKey, id string) (*gelatin.GelatinUser, error) {
	url := fmt.Sprintf("%s%s/%s", c.hostname, jellyfinUserGetEndpoint, id)
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

func (c *JellyfinApiClient) GetUsers(key gelatin.ApiKey, public bool) ([]gelatin.GelatinUser, error) {
	var url string
	if public {
		url = fmt.Sprintf("%s%s", c.hostname, jellyfinUserQueryPublicEndpoint)
	} else {
		url = fmt.Sprintf("%s%s", c.hostname, jellyfinUserQueryEndpoint)
	}

	raw, err := c.get(url, key)
	if err != nil {
		return nil, err
	}

	var resp []gelatin.GelatinUser
	dec := json.NewDecoder(raw.Body)
	if err := dec.Decode(&resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *JellyfinApiClient) UpdateUser(key gelatin.AdminKey, id string, data *gelatin.GelatinUser) error {
	url := fmt.Sprintf("%s%s/%s", c.hostname, jellyfinUserUpdateEndpoint, id)

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

func (c *JellyfinApiClient) NewUser(key gelatin.AdminKey, name string) (*gelatin.GelatinUser, error) {
	type createUserByName struct {
		Name string
	}
	req := createUserByName{Name: name}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s%s", c.hostname, jellyfinUserNewEndpoint)
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

func (c *JellyfinApiClient) DeleteUser(key gelatin.AdminKey, id string) error {
	url := fmt.Sprintf("%s%s/%s", c.hostname, jellyfinUserDeleteEndpoint, id)

	_, err := c.request(http.MethodDelete, url, nil, key)
	if err != nil {
		return err
	}

	return nil
}

func (c *JellyfinApiClient) UpdatePassword(key gelatin.AdminKey, id, currentPassword, newPassword string, reset bool) error {
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

	url := fmt.Sprintf("%s%s/%s/Password", c.hostname, jellyfinUserPasswordEndpoint, id)
	_, err = c.request(http.MethodPost, url, bytes.NewReader(data), key)
	if err != nil {
		return err
	}

	return nil
}

func (c *JellyfinApiClient) Authenticate(username, password string) (userKey gelatin.ApiKey, err error) {
	req := map[string]string{
		"Username": username,
		"Pw":       password,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s%s", c.hostname, jellyfinUserAuthEndpoint)
	raw, err := c.request(http.MethodPost, url, bytes.NewReader(data), nil)
	if err != nil {
		return nil, err
	}

	type authenticationResult struct {
		AccessToken string
	}

	resp := &authenticationResult{}

	dec := json.NewDecoder(raw.Body)
	if err := dec.Decode(resp); err != nil {
		return nil, err
	}

	return NewApiKey(resp.AccessToken), nil
}

func (c *JellyfinApiClient) UpdatePolicy(key gelatin.AdminKey, userId string, policy *gelatin.GelatinUserPolicy) error {
	url := fmt.Sprintf("%s%s/%s/Policy", c.hostname, jellyfinUserPolicyEndpoint, userId)

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
