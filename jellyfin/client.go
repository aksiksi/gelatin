package jellyfin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

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

const (
	jellyfinItemFilterParentId           = "parentId"
	jellyfinItemFilterFields             = "fields"
	jellyfinItemFilterRecursive          = "recursive"
	jellyfinItemFilterUserId             = "userId"
	jellyfinItemFilterFilters            = "filters"
	jellyfinItemFilterFiltersIsFolder    = "IsFolder"
	jellyfinItemFilterFiltersIsNotFolder = "IsNotFolder"
	jellyfinItemFilterFiltersIsPlayed    = "IsPlayed"

	jellyfinProviderIdImdb = "imdb"
	jellyfinProviderIdTmdb = "tmdb"
	jellyfinProviderIdTvdb = "tvdb"
)

type jellyfinApiKey struct {
	key     string
	isAdmin bool
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

func (k *jellyfinApiKey) IsAdmin() bool {
	return k.isAdmin
}

type JellyfinApiClient struct {
	client   *http.Client
	hostname string
	apiKey   gelatin.ApiKey
	mu       sync.Mutex
}

func NewJellyfinApiClient(hostname string, apiKey gelatin.ApiKey) *JellyfinApiClient {
	return &JellyfinApiClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		hostname: hostname,
		apiKey:   apiKey,
	}
}

func (c *JellyfinApiClient) ApiKey() gelatin.ApiKey {
	return c.apiKey
}

func (c *JellyfinApiClient) SetApiKey(key gelatin.ApiKey) {
	c.mu.Lock()
	c.apiKey = key
	c.mu.Unlock()
}

func (c *JellyfinApiClient) System() gelatin.GelatinSystemService {
	// TODO: Move this out
	return c
}

func (c *JellyfinApiClient) User() gelatin.GelatinUserService {
	// TODO: Move this out
	return c
}

func (c *JellyfinApiClient) Library() gelatin.GelatinLibraryService {
	// TODO: Move this out
	return c
}

func (c *JellyfinApiClient) Playlist() gelatin.GelatinPlaylistService {
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
	resp, err := c.Info(true)
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

func (c *JellyfinApiClient) GetLogs() ([]gelatin.GelatinSystemLog, error) {
	url := fmt.Sprintf("%s%s", c.hostname, jellyfinSystemLogsEndpoint)
	raw, err := c.get(url, c.apiKey)
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

func (c *JellyfinApiClient) GetLogFile(name string) (io.ReadCloser, error) {
	url := fmt.Sprintf("%s%s?name=%s", c.hostname, jellyfinSystemLogsNameEndpoint, name)

	resp, err := c.get(url, c.apiKey)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (c *JellyfinApiClient) Info(public bool) (*gelatin.GelatinSystemInfo, error) {
	var url string
	if public {
		url = fmt.Sprintf("%s%s", c.hostname, jellyfinSystemInfoPublicEndpoint)
	} else {
		url = fmt.Sprintf("%s%s", c.hostname, jellyfinSystemInfoEndpoint)
	}

	raw, err := c.get(url, c.apiKey)
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

func (c *JellyfinApiClient) GetUser(id string) (*gelatin.GelatinUser, error) {
	url := fmt.Sprintf("%s%s/%s", c.hostname, jellyfinUserGetEndpoint, id)
	raw, err := c.get(url, c.apiKey)
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

func (c *JellyfinApiClient) GetUsers(public bool) ([]gelatin.GelatinUser, error) {
	var url string
	if public {
		url = fmt.Sprintf("%s%s", c.hostname, jellyfinUserQueryPublicEndpoint)
	} else {
		url = fmt.Sprintf("%s%s", c.hostname, jellyfinUserQueryEndpoint)
	}

	raw, err := c.get(url, c.apiKey)
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

func (c *JellyfinApiClient) UpdateUser(id string, data *gelatin.GelatinUser) error {
	url := fmt.Sprintf("%s%s/%s", c.hostname, jellyfinUserUpdateEndpoint, id)

	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = c.request(http.MethodPost, url, bytes.NewReader(raw), c.apiKey)
	if err != nil {
		return err
	}

	return nil
}

func (c *JellyfinApiClient) CreateUser(name string) (*gelatin.GelatinUser, error) {
	type createUserByName struct {
		Name string
	}
	req := createUserByName{Name: name}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s%s", c.hostname, jellyfinUserNewEndpoint)
	raw, err := c.request(http.MethodPost, url, bytes.NewReader(data), c.apiKey)
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

func (c *JellyfinApiClient) DeleteUser(id string) error {
	url := fmt.Sprintf("%s%s/%s", c.hostname, jellyfinUserDeleteEndpoint, id)

	_, err := c.request(http.MethodDelete, url, nil, c.apiKey)
	if err != nil {
		return err
	}

	return nil
}

func (c *JellyfinApiClient) UpdatePassword(id, currentPassword, newPassword string, reset bool) error {
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
	_, err = c.request(http.MethodPost, url, bytes.NewReader(data), c.apiKey)
	if err != nil {
		return err
	}

	return nil
}

func (c *JellyfinApiClient) Authenticate(username, password string) (key gelatin.ApiKey, err error) {
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

	resp := &JellyfinUserAuthResponse{}

	dec := json.NewDecoder(raw.Body)
	if err := dec.Decode(resp); err != nil {
		return nil, err
	}

	return NewApiKey(resp.AccessToken), nil
}

func (c *JellyfinApiClient) UpdatePolicy(userId string, policy *gelatin.GelatinUserPolicy) error {
	url := fmt.Sprintf("%s%s/%s/Policy", c.hostname, jellyfinUserPolicyEndpoint, userId)

	data, err := json.Marshal(policy)
	if err != nil {
		return err
	}

	_, err = c.request(http.MethodPost, url, bytes.NewReader(data), c.apiKey)
	if err != nil {
		return err
	}

	return nil
}

func (c *JellyfinApiClient) GetItems(filters map[string]string, recursive bool) ([]gelatin.GelatinLibraryItem, error) {
	if filters == nil {
		filters = make(map[string]string)
	}

	endpoint := fmt.Sprintf("%s/Items", c.hostname)

	// Apply filters to URL string
	parsedUrl, _ := url.Parse(endpoint)
	query := parsedUrl.Query()
	for k, v := range filters {
		// Always include the ProviderIds in each returned library item
		if k == jellyfinItemFilterFields && !strings.Contains(v, "ProviderIds") {
			v += ", ProviderIds"
		}

		query.Set(k, v)
	}

	if _, ok := filters[jellyfinItemFilterFields]; !ok {
		query.Set(jellyfinItemFilterFields, "ProviderIds")
	}
	query.Set(jellyfinItemFilterRecursive, "true")

	parsedUrl.RawQuery = query.Encode()

	raw, err := c.get(parsedUrl.String(), c.apiKey)
	if err != nil {
		return nil, err
	}

	resp := &JellyfinLibraryItemResponse{}
	dec := json.NewDecoder(raw.Body)
	if err := dec.Decode(resp); err != nil {
		return nil, err
	}

	// Include the provider IDs as struct fields
	for i := range resp.Items {
		item := &resp.Items[i]

		for provider, id := range item.ProviderIds {
			switch strings.ToLower(provider) {
			case jellyfinProviderIdImdb:
				item.ImdbId = id
			case jellyfinProviderIdTmdb:
				item.TmdbId = id
			case jellyfinProviderIdTvdb:
				item.TvdbId = id
			}
		}
	}

	return resp.Items, nil
}

func (c *JellyfinApiClient) GetItemsByUser(id string, filters map[string]string) ([]gelatin.GelatinLibraryItem, error) {
	if filters == nil {
		filters = make(map[string]string)
	}

	filters[jellyfinItemFilterUserId] = id

	return c.GetItems(filters, true)
}

func (c *JellyfinApiClient) UpdateItem(itemId string, item *gelatin.GelatinLibraryItem) error {
	url := fmt.Sprintf("%s/Items/%s", c.hostname, itemId)

	data, err := json.Marshal(item)
	if err != nil {
		return err
	}

	_, err = c.request(http.MethodPost, url, bytes.NewReader(data), nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *JellyfinApiClient) updateItemFavoriteState(itemId string, userId string, favorite bool) error {
	url := fmt.Sprintf("%s/Users/%s/FavoriteItems/%s", c.hostname, userId, itemId)

	var method string
	if favorite {
		method = http.MethodPost
	} else {
		method = http.MethodDelete
	}

	_, err := c.request(method, url, nil, c.apiKey)
	if err != nil {
		return err
	}

	return nil
}

func (c *JellyfinApiClient) updateItemPlayedState(itemId string, userId string, played bool) error {
	url := fmt.Sprintf("%s/Users/%s/PlayedItems/%s", c.hostname, userId, itemId)

	var method string
	if played {
		method = http.MethodPost
	} else {
		method = http.MethodDelete
	}

	_, err := c.request(method, url, nil, c.apiKey)
	if err != nil {
		return err
	}

	return nil
}

func (c *JellyfinApiClient) updateItemPlayingState(itemId string, userId string, ticks int64) error {
	url := fmt.Sprintf("%s/Users/%s/PlayingItems/%s/Progress", c.hostname, userId, itemId)

	playingStateRequest := map[string]string{
		"positionTicks": strconv.FormatInt(ticks, 10),
	}

	data, err := json.Marshal(playingStateRequest)
	if err != nil {
		return err
	}

	_, err = c.request(http.MethodPost, url, bytes.NewReader(data), c.apiKey)
	if err != nil {
		return err
	}

	return nil
}

func (c *JellyfinApiClient) UpdateItemUserActivity(itemId string, userId string, old, new *gelatin.GelatinLibraryItemUserActivity) error {
	/*
		NOTE(aksiksi): Jellyfin does not expose a UserData update endpoint. So, to achieve the same thing,
		we need to use 3 different endpoints.

		- Favorite: /Users/{userId}/PlayedItems/{itemId} (POST)
			- DELETE to remove
		- Played: /Users/{userId}/FavoriteItems/{itemId} (POST)
			- DELETE to remove

		As far as watch progress goes, you need to start a play session and report an update...
	*/
	if old.IsFavorite != new.IsFavorite {
		if err := c.updateItemFavoriteState(itemId, userId, new.IsFavorite); err != nil {
			return err
		}
	}

	if old.Played != new.Played {
		if err := c.updateItemPlayedState(itemId, userId, new.Played); err != nil {
			return err
		}
	}

	// TODO(aksiksi): Figure out why this isn't working. Do we need to use /Sessions?
	if !new.Played && old.PlaybackPositionTicks != new.PlaybackPositionTicks {
		if err := c.updateItemPlayingState(itemId, userId, new.PlaybackPositionTicks); err != nil {
			return err
		}
	}

	return nil
}

func (c *JellyfinApiClient) GetItemFilterString(filter gelatin.GelatinItemFilterName) string {
	switch filter {
	case gelatin.GelatinItemFilterFilters:
		return jellyfinItemFilterFilters
	case gelatin.GelatinItemFilterParentId:
		return jellyfinItemFilterParentId
	case gelatin.GelatinItemFilterFiltersIsFolder:
		return jellyfinItemFilterFiltersIsFolder
	case gelatin.GelatinItemFilterFiltersIsNotFolder:
		return jellyfinItemFilterFiltersIsNotFolder
	case gelatin.GelatinItemFilterFiltersIsPlayed:
		return jellyfinItemFilterFiltersIsPlayed
	default:
		panic("invalid filter name")
	}
}
