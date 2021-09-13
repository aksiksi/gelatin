package gelatin

import (
	"io"
)

type GelatinSystemLog struct {
	Name         string
	Size         int64
	DateCreated  string
	DateModified string
}

type GelatinSystemInfo struct {
	Id              string
	LocalAddress    string
	WanAddress      string
	ServerName      string
	Version         string
	OperatingSystem string
}

type GelatinUserAccessSchedule struct {
	DayOfWeek string `validate:"oneof=Sunday Monday Tuesday Wednesday Thursday Friday Saturday Everyday Weekday Weekend"`
	StartHour float64
	EndHour   float64
}

type GelatinUserPolicy struct {
	IsAdministrator                  bool
	IsHidden                         bool
	IsHiddenRemotely                 bool
	IsHiddenFromUnusedDevices        bool
	IsDisabled                       bool
	MaxParentalRating                *int32
	BlockedTags                      []string
	IsTagBlockingModeInclusive       bool
	EnableUserPreferenceAccess       bool
	AccessSchedules                  []GelatinUserAccessSchedule
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

type GelatinUserConfig struct {
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

type GelatinUser struct {
	Name                      string
	ServerId                  string
	ServerName                string
	Id                        string `validate:"uuid"`
	PrimaryImageTag           string
	HasPassword               bool
	HasConfiguredPassword     bool
	HasConfiguredEasyPassword bool
	EnableAutoLogin           *bool
	LastLoginDate             *string `validate:"datetime=2006-01-02T15:04:05.0000000Z"`
	LastActivityDate          *string `validate:"datetime=2006-01-02T15:04:05.0000000Z"`
	Configuration             GelatinUserConfig
	Policy                    GelatinUserPolicy
	PrimaryImageAspectRatio   *float64
}

type GelatinSystemService interface {
	Version() (string, error)
	Ping() error
	GetLogs(key ApiKey) ([]GelatinSystemLog, error)
	GetLogFile(key ApiKey, name string) (io.ReadCloser, error)
	Info(key ApiKey, public bool) (*GelatinSystemInfo, error)
}

type GelatinUserService interface {
	GetUser(key ApiKey, id string) (*GelatinUser, error)
	GetUsers(key ApiKey, public bool) ([]GelatinUser, error)
	UpdateUser(key AdminKey, id string, data *GelatinUser) error
	NewUser(key AdminKey, name string) (*GelatinUser, error)
	DeleteUser(key AdminKey, id string) error
	UpdatePassword(key AdminKey, id, currentPassword, newPassword string, reset bool) error
	Authenticate(username, password string) (userKey ApiKey, err error)
	UpdatePolicy(key AdminKey, id string, policy *GelatinUserPolicy) error
}

type GelatinLibraryService interface {
}

type GelatinPlaylistService interface {
}

type GelatinService interface {
	System() GelatinSystemService
	User() GelatinUserService
	Library() GelatinLibraryService
	Playlist() GelatinPlaylistService
}
