package jellyfin

import (
	"io"

	"github.com/aksiksi/gelatin/api"
)

const (
	jellyfinApiKeyHeaderName = "X-Jellyfin-Token"
)

const (
	jellyfinSystemPingEndpoint       = "/System/Ping"
	jellyfinSystemLogsEndpoint       = "/System/Logs"
	jellyfinSystemLogsQueryEndpoint  = "/System/Logs/Query"
	jellyfinSystemInfoEndpoint       = "/System/Info"
	jellyfinSystemInfoPublicEndpoint = "/System/Info/Public"
)

type JellyfinSystemServiceApi interface {
	SystemPing() error
	SystemLogs(key api.ApiKey, name string) (io.ReadCloser, error)
	SystemLogsQuery(key api.ApiKey) (*JellyfinSystemLogsQueryResponse, error)
	SystemInfo(key api.ApiKey) (*JellyfinSystemInfoResponse, error)
	SystemInfoPublic() (*JellyfinSystemInfoPublicResponse, error)
}

type JellyfinApi interface {
	JellyfinSystemServiceApi
}
