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
	PasswordResetProviderId          string // Jellyfin only
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

// GelatinLibraryItemUserActivity contains user activity associated with a
// given library item.
type GelatinLibraryItemUserActivity struct {
	PlaybackPositionTicks int64
	PlayCount             int32
	IsFavorite            bool
	LastPlayedDate        string
	Played                bool
	Rating                float64
	UnplayedItemCount     int32 // For series
}

// GelatinLibraryItem holds info for a single library item
//
// Note that this struct tracks a subset of available fields from Emby and Jellyfin.
// This is by design: we are only interested in getting and updating a few fields for
// a given library item.
type GelatinLibraryItem struct {
	Name         string
	ServerId     string
	Id           string
	RunTimeTicks int64
	IsFolder     bool
	Type         string // Movie, Series, Episode, etc.
	UserData     *GelatinLibraryItemUserActivity
	MediaType    string // Video, Photo, etc.

	ProviderIds map[string]string
	ImdbId      string
	TmdbId      string
	TvdbId      string

	// Series-specific info (Type == "Series")
	SeriesId          string // ID of the series this episode belongs to
	SeriesName        string
	SeasonId          string
	EpisodeNumber     int32
	IndexNumber       int32
	ParentIndexNumber int32
}

type GelatinSystemService interface {
	// Version returns the version string
	Version() (string, error)

	// Ping pings the server (i.e., health check)
	Ping() error

	// GetLogs returns all logs exposed by the server
	GetLogs() ([]GelatinSystemLog, error)

	// GetLogFile downloads the content of a single log file
	GetLogFile(name string) (io.ReadCloser, error)

	// Info returns information about the server
	//
	// If "public" is true, this returns only publicly visible system info.
	Info(public bool) (*GelatinSystemInfo, error)
}

type GelatinUserService interface {
	// GetUser returns the user with the specified ID
	GetUser(id string) (*GelatinUser, error)

	// GetUsers returns all configured users
	//
	// If "public" is true, returns only the publicly visible users.
	GetUsers(public bool) ([]GelatinUser, error)

	// UpdateUser updates a single user
	//
	// Note that user state is _overwritten_. Use this in conjunction with
	// GetUser().
	UpdateUser(id string, data *GelatinUser) error

	// CreateUser creates a new user with the given username
	CreateUser(name string) (*GelatinUser, error)

	// DeleteUser deletes the user with the specified ID
	DeleteUser(id string) error

	// UpdatePassword updates the given user's password
	//
	// If "reset" is true, the password will be reset first.
	UpdatePassword(id, currentPassword, newPassword string, reset bool) error

	// Authenticate as a specific user
	//
	// Use this method with an admin account to create an AdminKey.
	Authenticate(username, password string) (key ApiKey, err error)

	// UpdatePolicy updates the policy for the specified user.
	//
	// Note that user state is _overwritten_. Use this in conjunction with
	// GetUser().
	UpdatePolicy(id string, policy *GelatinUserPolicy) error
}

type GelatinLibraryService interface {
	// GetItems returns library items for a _specific_ user (i.e., with user activity attached)
	//
	// Refer to Emby or Jellyfin docs for available item filters. Note that the "recursive"
	// filter will always be present.
	GetItems(id string, filters map[string]string) ([]GelatinLibraryItem, error)
}

type GelatinPlaylistService interface {
}

type GelatinService interface {
	// ApiKey returns the current API key used by the client
	ApiKey() ApiKey
	SetApiKey(key ApiKey)

	System() GelatinSystemService
	User() GelatinUserService
	Library() GelatinLibraryService
	Playlist() GelatinPlaylistService
}
