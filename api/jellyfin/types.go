package jellyfin

type JellyfinInstallationInfo struct {
	Id              string `validate:"uuid"`
	Name            string
	AssemblyGuid    string
	Version         string
	UpdateClass     string `validate:"oneof=Release Beta Dev"`
	PercentComplete *float64
}

type JellyfinSystemInfoResponse struct {
	SystemUpdateLevel                    string `validate:"oneof=Release Beta Dev"`
	OperatingSystemDisplayName           string
	PackageName                          string
	HasPendingRestart                    bool
	IsShuttingDown                       bool
	SupportsLibraryMonitor               bool
	WebSocketPortNumber                  int32
	CompletedInstallations               JellyfinInstallationInfo
	CanSelfRestart                       bool
	CanSelfUpdate                        bool
	CanLaunchWebBrowser                  bool
	ProgramDataPath                      string
	ItemsByNamePath                      string
	CachePath                            string
	LogPath                              string
	InternalMetadataPath                 string
	TranscodingTempPath                  string
	HttpServerPortNumber                 int32
	SupportsHttps                        bool
	HttpsPortNumber                      int32
	HasUpdateAvailable                   bool
	SupportsAutoRunAtStartup             bool
	HardwareAccelerationRequiresPremiere bool
	JellyfinSystemInfoPublicResponse
}

type JellyfinSystemInfoPublicResponse struct {
	LocalAddress    string
	WanAddress      string
	ServerName      string
	Version         string
	OperatingSystem string
	Id              string
}
