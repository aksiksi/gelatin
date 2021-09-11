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
	jellyfinUserUpdateEndpoint       = "/users"
	jellyfinUserNewEndpoint          = "/users/new"
	jellyfinUserDeleteEndpoint       = "/users"
	jellyfinUserPasswordEndpoint     = "/Users"
	jellyfinUserAuthEndpoint         = "/Users/AuthenticateByName"
	jellyfinUserPolicyEndpoint       = "/Users"
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
	UserUpdate(key api.AdminKey, userId string, dto *JellyfinUserDto) error
	UserNew(key api.AdminKey, name string) (*JellyfinUserDto, error)
	UserDelete(key api.AdminKey, userId string) error
	UserPassword(key api.AdminKey, userId, currentPassword, newPassword string, reset bool) error
	UserAuth(username, password string) (userKey api.ApiKey, err error)
	UserPolicy(key api.AdminKey, userId string, policy *JellyfinUserPolicy) error
}

type JellyfinApi interface {
	JellyfinSystemServiceApi
	JellyfinUserServiceApi
}
