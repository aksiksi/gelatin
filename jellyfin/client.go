package jellyfin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

func (c *JellyfinApiClient) GetVersion() (string, error) {
	resp, err := c.SystemInfoPublic()
	if err != nil {
		return "", err
	}

	return resp.Version, nil
}

func (c *JellyfinApiClient) request(method string, url string, body io.Reader, key gelatin.ApiKey) (*http.Response, error) {
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

	if err := gelatin.HttpStatusToErr(resp.StatusCode); err != nil {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("response body: %s", body)
		return nil, err
	}

	return resp, nil
}

func (c *JellyfinApiClient) get(url string, key gelatin.ApiKey) (*http.Response, error) {
	resp, err := c.request(http.MethodGet, url, nil, key)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *JellyfinApiClient) SystemPing() error {
	url := fmt.Sprintf("%s%s", c.hostname, jellyfinSystemPingEndpoint)
	_, err := c.get(url, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *JellyfinApiClient) SystemLogsName(key gelatin.ApiKey, name string) (io.ReadCloser, error) {
	url := fmt.Sprintf("%s%s?name=%s", c.hostname, jellyfinSystemLogsNameEndpoint, name)

	resp, err := c.get(url, key)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (c *JellyfinApiClient) SystemLogs(key gelatin.ApiKey) ([]JellyfinSystemLogFile, error) {
	url := fmt.Sprintf("%s%s", c.hostname, jellyfinSystemLogsEndpoint)
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

func (c *JellyfinApiClient) SystemInfo(key gelatin.ApiKey) (*JellyfinSystemInfoResponse, error) {
	url := fmt.Sprintf("%s%s", c.hostname, jellyfinSystemInfoEndpoint)
	raw, err := c.get(url, key)
	if err != nil {
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
	url := fmt.Sprintf("%s%s", c.hostname, jellyfinUserQueryPublicEndpoint)
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

func (c *JellyfinApiClient) UserQuery(key gelatin.ApiKey) ([]*JellyfinUserDto, error) {
	url := fmt.Sprintf("%s%s", c.hostname, jellyfinUserQueryEndpoint)
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

func (c *JellyfinApiClient) UserGet(key gelatin.ApiKey, userId string) (*JellyfinUserDto, error) {
	url := fmt.Sprintf("%s%s/%s", c.hostname, jellyfinUserGetEndpoint, userId)
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

func (c *JellyfinApiClient) UserUpdate(key gelatin.AdminKey, userId string, dto *JellyfinUserDto) error {
	url := fmt.Sprintf("%s%s/%s", c.hostname, jellyfinUserUpdateEndpoint, userId)

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

func (c *JellyfinApiClient) UserNew(key gelatin.ApiKey, name string) (*JellyfinUserDto, error) {
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

	resp := &JellyfinUserDto{}
	dec := json.NewDecoder(raw.Body)
	if err := dec.Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *JellyfinApiClient) UserDelete(key gelatin.ApiKey, userId string) error {
	url := fmt.Sprintf("%s%s/%s", c.hostname, jellyfinUserDeleteEndpoint, userId)

	_, err := c.request(http.MethodDelete, url, nil, key)
	if err != nil {
		return err
	}

	return nil
}

func (c *JellyfinApiClient) ResetUserPassword(key gelatin.AdminKey, userId string) error {
	type resetUserPassword struct {
		Id            string
		ResetPassword bool
	}

	req := resetUserPassword{Id: userId, ResetPassword: true}
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s%s/%s/Password", c.hostname, jellyfinUserPasswordEndpoint, userId)
	_, err = c.request(http.MethodPost, url, bytes.NewReader(data), key)
	if err != nil {
		return err
	}

	return nil
}

func (c *JellyfinApiClient) UserPassword(key gelatin.AdminKey, userId, currentPassword, newPassword string, reset bool) error {
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

	url := fmt.Sprintf("%s%s/%s/Password", c.hostname, jellyfinUserPasswordEndpoint, userId)
	_, err = c.request(http.MethodPost, url, bytes.NewReader(data), key)
	if err != nil {
		return err
	}

	return nil
}

func (c *JellyfinApiClient) UserAuth(username, password string) (userKey gelatin.ApiKey, err error) {
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

func (c *JellyfinApiClient) UserPolicy(key gelatin.AdminKey, userId string, policy *JellyfinUserPolicy) error {
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
