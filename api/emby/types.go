package emby

type EmbySystemLogFile struct {
	DateCreated  string `validate:"datetime=2006-01-02T15:04:05.0000000Z"`
	DateModified string `validate:"datetime=2006-01-02T15:04:05.0000000Z"`
	Size         int64
	Name         string
}

type EmbySystemLogsQueryResponse struct {
	Items            []EmbySystemLogFile
	TotalRecordCount int32
}

type EmbyInstallationInfo struct {
	Id              string `validate:"uuid"`
	Name            string
	AssemblyGuid    string
	Version         string
	UpdateClass     string `validate:"oneof=Release Beta Dev"`
	PercentComplete *float64
}

type EmbySystemInfoResponse struct {
	SystemUpdateLevel                    string `validate:"oneof=Release Beta Dev"`
	OperatingSystemDisplayName           string
	PackageName                          string
	HasPendingRestart                    bool
	IsShuttingDown                       bool
	SupportsLibraryMonitor               bool
	WebSocketPortNumber                  int32
	CompletedInstallations               EmbyInstallationInfo
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
	EmbySystemInfoPublicResponse
}

type EmbySystemInfoPublicResponse struct {
	LocalAddress    string
	WanAddress      string
	ServerName      string
	Version         string
	OperatingSystem string
	Id              string
}
