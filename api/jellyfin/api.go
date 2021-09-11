package jellyfin

import (
	"io"

	"github.com/aksiksi/gelatin/api"
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
	jellyfinUserUpdateEndpoint       = jellyfinUserGetEndpoint
	jellyfinUserNewEndpoint          = "/users/new"
	jellyfinUserPasswordEndpoint     = "/Users"
)

type JellyfinSystemServiceApi interface {
	SystemPing() error
	SystemLogs(key api.ApiKey) ([]JellyfinSystemLogFile, error)
	SystemLogsName(key api.ApiKey, name string) (io.ReadCloser, error)
	SystemInfo(key api.ApiKey) (*JellyfinSystemInfoResponse, error)
	SystemInfoPublic() (*JellyfinSystemInfoPublicResponse, error)
	GetVersion() (string, error)
}

type JellyfinUserServiceApi interface {
	UserQueryPublic() ([]*JellyfinUserDto, error)
	UserQuery(key api.ApiKey) ([]*JellyfinUserDto, error)
	UserGet(key api.ApiKey, userId string) (*JellyfinUserDto, error)
	UserUpdate(key api.ApiKey, userId string, dto *JellyfinUserDto) error
	UserNew(key api.ApiKey, name string) (*JellyfinUserDto, error)
	UserPassword(key api.ApiKey, userId, currentPassword, newPassword string, reset bool) error
}

type JellyfinApi interface {
	JellyfinSystemServiceApi
	JellyfinUserServiceApi
}
