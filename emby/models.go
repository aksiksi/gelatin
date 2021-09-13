package emby

import gelatin "github.com/aksiksi/gelatin/lib"

type EmbyUserAuthResponse struct {
	AccessToken string
	// Other fields ommitted
}

type EmbyUserAccessSchedule struct {
	DayOfWeek string `validate:"oneof=Sunday Monday Tuesday Wednesday Thursday Friday Saturday Everyday Weekday Weekend"`
	StartHour float64
	EndHour   float64
}

type EmbyUserPolicy struct {
	IsAdministrator                  bool
	IsHidden                         bool
	IsHiddenRemotely                 bool
	IsHiddenFromUnusedDevices        bool
	IsDisabled                       bool
	MaxParentalRating                *int32
	BlockedTags                      []string
	IsTagBlockingModeInclusive       bool
	EnableUserPreferenceAccess       bool
	AccessSchedules                  []EmbyUserAccessSchedule
	BlockUnratedItems                []string
	EnableRemoteControlOfOtherUsers  bool
	EnableSharedDeviceControl        bool
	EnableRemoteAccess               bool
	EnableLiveTvManagement           bool
	EnableLiveTvAccess               bool
	EnableMediaPlayback              bool
	EnableAudioPlaybackTranscoding   bool
	EnableVideoPlaybackTranscoding   bool
	EnablePlaybackRemuxing           bool
	EnableContentDeletion            bool
	EnableContentDeletionFromFolders []string
	EnableContentDownloading         bool
	EnableSubtitleDownloading        bool
	EnableSubtitleManagement         bool
	EnableSyncTranscoding            bool
	EnableMediaConversion            bool
	EnabledChannels                  []string
	EnableAllChannels                bool
	EnabledFolders                   []string
	EnableAllFolders                 bool
	InvalidLoginAttemptCount         int32
	EnablePublicSharing              bool
	BlockedMediaFolders              []string
	BlockedChannels                  []string
	RemoteClientBitrateLimit         int32
	AuthenticationProviderId         string
	ExcludedSubFolders               []string
	SimultaneousStreamLimit          int32
	EnabledDevices                   []string
	EnableAllDevices                 bool
}

type EmbyUserConfig struct {
	AudioLanguagePreference    string
	PlayDefaultAudioTrack      bool
	SubtitleLanguagePreference string
	DisplayMissingEpisodes     bool
	SubtitleMode               string `validate:"oneof=Default Always OnlyForced None Smart"`
	EnableLocalPassword        bool
	OrderedViews               []string
	LatestItemsExcludes        []string
	MyMediaExcludes            []string
	HidePlayedInLatest         bool
	RememberAudioSelections    bool
	RememberSubtitleSelections bool
	EnableNextEpisodeAutoPlay  bool
}

type EmbyUserDto struct {
	Name                      string
	ServerId                  string
	ServerName                string
	ConnectUserName           string
	ConnectLinkType           string `validate:"oneof=LinkedUser Guest"`
	Id                        string `validate:"uuid"`
	PrimaryImageTag           string
	HasPassword               bool
	HasConfiguredPassword     bool
	HasConfiguredEasyPassword bool
	EnableAutoLogin           *bool
	LastLoginDate             *string `validate:"datetime=2006-01-02T15:04:05.0000000Z"`
	LastActivityDate          *string `validate:"datetime=2006-01-02T15:04:05.0000000Z"`
	Configuration             EmbyUserConfig
	Policy                    EmbyUserPolicy
	PrimaryImageAspectRatio   *float64
}

type EmbyUserQueryResponse struct {
	Items            []gelatin.GelatinUser
	TotalRecordCount int32
}

type EmbySystemLogFile struct {
	DateCreated  string `validate:"datetime=2006-01-02T15:04:05.0000000Z"`
	DateModified string `validate:"datetime=2006-01-02T15:04:05.0000000Z"`
	Size         int64
	Name         string
}

type EmbySystemLogsQueryResponse struct {
	Items            []gelatin.GelatinSystemLog
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
