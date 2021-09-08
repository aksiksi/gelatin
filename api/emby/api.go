package emby

import (
	"io"

	"github.com/aksiksi/gelatin/api"
)

const (
	embyApiKeyHeaderName = "X-Emby-Token"
)

const (
	embySystemPingEndpoint       = "/System/Ping"
	embySystemLogsEndpoint       = "/System/Logs"
	embySystemLogsQueryEndpoint  = "/System/Logs/Query"
	embySystemInfoEndpoint       = "/System/Info"
	embySystemInfoPublicEndpoint = "/System/Info/Public"
)

type EmbySystemServiceApi interface {
	SystemPing() error
	SystemLogs(key api.ApiKey, name string) (io.ReadCloser, error)
	SystemLogsQuery(key api.ApiKey) (*EmbySystemLogsQueryResponse, error)
	SystemInfo(key api.ApiKey) (*EmbySystemInfoResponse, error)
	SystemInfoPublic() (*EmbySystemInfoPublicResponse, error)
}

type EmbyApi interface {
	EmbySystemServiceApi
}
