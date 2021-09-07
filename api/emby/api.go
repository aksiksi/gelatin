package emby

const (
	embySystemPingEndpoint       = "/System/Ping"
	embySystemLogsEndpoint       = "/System/Logs"
	embySystemInfoEndpoint       = "/System/Info"
	embySystemInfoPublicEndpoint = "/System/Info/Public"
)

type EmbySystemServiceApi interface {
	SystemPing() error
	SystemLogs() error
	SystemInfo() error
	SystemInfoPublic() error
}

type EmbyApi interface {
	EmbySystemServiceApi
}
