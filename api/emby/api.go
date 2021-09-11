package emby

import (
	"io"

	"github.com/aksiksi/gelatin/api"
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
	embyUserQueryPublicEndpoint  = "/Users/Query/Public"
	embyUserGetEndpoint          = "/Users"
	embyUserUpdateEndpoint       = embyUserGetEndpoint
	embyUserNewEndpoint          = "/Users/New"
)

type EmbySystemServiceApi interface {
	SystemPing() error
	SystemLogs(key api.ApiKey, name string) (io.ReadCloser, error)
	SystemLogsQuery(key api.ApiKey) (*EmbySystemLogsQueryResponse, error)
	SystemInfo(key api.ApiKey) (*EmbySystemInfoResponse, error)
	SystemInfoPublic() (*EmbySystemInfoPublicResponse, error)
	GetVersion() (string, error)
}

type EmbyUserServiceApi interface {
	UserQueryPublic() (*EmbyUserQueryResponse, error)
	UserQuery(key api.ApiKey) (*EmbyUserQueryResponse, error)
	UserGet(key api.ApiKey, userId string) (*EmbyUserDto, error)
	UserUpdate(key api.ApiKey, userId string, dto *EmbyUserDto) error
	UserNew(key api.ApiKey, name string) (*EmbyUserDto, error)
}

type EmbyApi interface {
	EmbySystemServiceApi
	EmbyUserServiceApi
}
