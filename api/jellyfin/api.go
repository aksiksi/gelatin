package jellyfin

const (
	jellyfinSystemPingEndpoint       = "/System/Ping"
	jellyfinSystemLogsEndpoint       = "/System/Logs"
	jellyfinSystemInfoEndpoint       = "/System/Info"
	jellyfinSystemInfoPublicEndpoint = "/System/Info/Public"
)

type JellyfinSystemServiceApi interface {
	SystemPing() error
	SystemLogs() error
	SystemInfo() error
	SystemInfoPublic() error
}

type JellyfinApi interface {
	JellyfinSystemServiceApi
}
