package jellyfin

import (
	"io"

	"github.com/aksiksi/gelatin/api"
)

const (
	// See: https://github.com/jellyfin/jellyfin/blob/8c463b9b8196530e6d417b40ed81825422cada32/Jellyfin.Server.Implementations/Security/AuthorizationContext.cs#L80
	jellyfinApiKeyHeaderName = "X-Emby-Token"
)

const (
	jellyfinSystemPingEndpoint       = "/System/Ping"
	jellyfinSystemLogsEndpoint       = "/System/Logs"
	jellyfinSystemLogsNameEndpoint   = "/System/Logs/Log"
	jellyfinSystemInfoEndpoint       = "/System/Info"
	jellyfinSystemInfoPublicEndpoint = "/System/Info/Public"
)

type JellyfinSystemServiceApi interface {
	SystemPing() error
	SystemLogs(key api.ApiKey) ([]JellyfinSystemLogFile, error)
	SystemLogsName(key api.ApiKey, name string) (io.ReadCloser, error)
	SystemInfo(key api.ApiKey) (*JellyfinSystemInfoResponse, error)
	SystemInfoPublic() (*JellyfinSystemInfoPublicResponse, error)
}

type JellyfinApi interface {
	JellyfinSystemServiceApi
}
