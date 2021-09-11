package jellyfin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/aksiksi/gelatin/api"
)

type jellyfinApiKey struct {
	key string
}

// NewApiKey returns a new ApiKey for the given client
func NewApiKey(key string) api.ApiKey {
	return &jellyfinApiKey{
		key: key,
	}
}

func (k *jellyfinApiKey) ToString() string {
	return k.key
}

type JellyfinApiClient struct {
	client   *http.Client
	endpoint string
}

func NewJellyfinApiClient(hostname string, client *http.Client) *JellyfinApiClient {
	return &JellyfinApiClient{
		client:   client,
		endpoint: hostname,
	}
}

func (c *JellyfinApiClient) GetVersion() (string, error) {
	resp, err := c.SystemInfoPublic()
	if err != nil {
		return "", err
	}

	return resp.Version, nil
}

func (c *JellyfinApiClient) request(method string, url string, body io.Reader, key api.ApiKey) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add(jellyfinApiKeyAuthHeader, `MediaBrowser Client="gelatin", Device="gelatin", DeviceId="007", Version="0.0.1"`)

	if key != nil {
		req.Header.Add(jellyfinApiKeyTokenHeader, key.ToString())
	}

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if err := api.HttpStatusToErr(resp.StatusCode); err != nil {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("response body: %s", body)
		return nil, err
	}

	return resp, nil
}

func (c *JellyfinApiClient) get(url string, key api.ApiKey) (*http.Response, error) {
	resp, err := c.request(http.MethodGet, url, nil, key)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *JellyfinApiClient) SystemPing() error {
	url := fmt.Sprintf("%s%s", c.endpoint, jellyfinSystemPingEndpoint)
	_, err := c.get(url, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *JellyfinApiClient) SystemLogsName(key api.ApiKey, name string) (io.ReadCloser, error) {
	url := fmt.Sprintf("%s%s?name=%s", c.endpoint, jellyfinSystemLogsNameEndpoint, name)

	resp, err := c.get(url, key)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (c *JellyfinApiClient) SystemLogs(key api.ApiKey) ([]JellyfinSystemLogFile, error) {
	url := fmt.Sprintf("%s%s", c.endpoint, jellyfinSystemLogsEndpoint)
	raw, err := c.get(url, key)
	if err != nil {
		return nil, err
	}

	var resp []JellyfinSystemLogFile
	dec := json.NewDecoder(raw.Body)
	if err := dec.Decode(&resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *JellyfinApiClient) SystemInfo(key api.ApiKey) (*JellyfinSystemInfoResponse, error) {
	url := fmt.Sprintf("%s%s", c.endpoint, jellyfinSystemInfoEndpoint)
	raw, err := c.get(url, key)
	if err != nil {
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
	url := fmt.Sprintf("%s%s", c.endpoint, jellyfinSystemInfoPublicEndpoint)
	raw, err := c.get(url, nil)
	if err != nil {
		return nil, err
	}

	resp := &JellyfinSystemInfoPublicResponse{}
	dec := json.NewDecoder(raw.Body)
	if err := dec.Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *JellyfinApiClient) UserQueryPublic() ([]*JellyfinUserDto, error) {
	url := fmt.Sprintf("%s%s", c.endpoint, jellyfinUserQueryPublicEndpoint)
	raw, err := c.get(url, nil)
	if err != nil {
		return nil, err
	}

	var resp []*JellyfinUserDto
	dec := json.NewDecoder(raw.Body)
	if err := dec.Decode(&resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *JellyfinApiClient) UserQuery(key api.ApiKey) ([]*JellyfinUserDto, error) {
	url := fmt.Sprintf("%s%s", c.endpoint, jellyfinUserQueryEndpoint)
	raw, err := c.get(url, key)
	if err != nil {
		return nil, err
	}

	var resp []*JellyfinUserDto
	dec := json.NewDecoder(raw.Body)
	if err := dec.Decode(&resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *JellyfinApiClient) UserGet(key api.ApiKey, userId string) (*JellyfinUserDto, error) {
	url := fmt.Sprintf("%s%s/%s", c.endpoint, jellyfinUserGetEndpoint, userId)
	raw, err := c.get(url, key)
	if err != nil {
		return nil, err
	}

	resp := &JellyfinUserDto{}
	dec := json.NewDecoder(raw.Body)
	if err := dec.Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *JellyfinApiClient) UserUpdate(key api.AdminKey, userId string, dto *JellyfinUserDto) error {
	url := fmt.Sprintf("%s%s/%s", c.endpoint, jellyfinUserUpdateEndpoint, userId)

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

func (c *JellyfinApiClient) UserNew(key api.ApiKey, name string) (*JellyfinUserDto, error) {
	type createUserByName struct {
		Name string
	}
	req := createUserByName{Name: name}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s%s", c.endpoint, jellyfinUserNewEndpoint)
	raw, err := c.request(http.MethodPost, url, bytes.NewReader(data), key)
	if err != nil {
		return nil, err
	}

	resp := &JellyfinUserDto{}
	dec := json.NewDecoder(raw.Body)
	if err := dec.Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *JellyfinApiClient) UserDelete(key api.ApiKey, userId string) error {
	url := fmt.Sprintf("%s%s/%s", c.endpoint, jellyfinUserDeleteEndpoint, userId)

	_, err := c.request(http.MethodDelete, url, nil, key)
	if err != nil {
		return err
	}

	return nil
}

func (c *JellyfinApiClient) ResetUserPassword(key api.AdminKey, userId string) error {
	type resetUserPassword struct {
		Id            string
		ResetPassword bool
	}

	req := resetUserPassword{Id: userId, ResetPassword: true}
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s%s/%s/Password", c.endpoint, jellyfinUserPasswordEndpoint, userId)
	_, err = c.request(http.MethodPost, url, bytes.NewReader(data), key)
	if err != nil {
		return err
	}

	return nil
}

func (c *JellyfinApiClient) UserPassword(key api.AdminKey, userId, currentPassword, newPassword string, reset bool) error {
	type setUserPassword struct {
		Id        string
		CurrentPw string
		NewPw     string
	}

	// Apparently, you need to always reset the user's password, even if you're
	// just changing it...
	//
	// See this issue for details: https://github.com/jellyfin/jellyfin/issues/1297
	if err := c.ResetUserPassword(key, userId); err != nil {
		return err
	}

	log.Print("done")

	if reset {
		return nil
	}

	req := setUserPassword{
		Id:        userId,
		CurrentPw: currentPassword,
		NewPw:     newPassword,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s%s/%s/Password", c.endpoint, jellyfinUserPasswordEndpoint, userId)
	_, err = c.request(http.MethodPost, url, bytes.NewReader(data), key)
	if err != nil {
		return err
	}

	return nil
}

func (c *JellyfinApiClient) UserAuth(username, password string) (userKey api.ApiKey, err error) {
	req := map[string]string{
		"Username": username,
		"Pw":       password,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s%s", c.endpoint, jellyfinUserAuthEndpoint)
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

func (c *JellyfinApiClient) UserPolicy(key api.AdminKey, userId string, policy *JellyfinUserPolicy) error {
	url := fmt.Sprintf("%s%s/%s/Policy", c.endpoint, jellyfinUserPolicyEndpoint, userId)

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
