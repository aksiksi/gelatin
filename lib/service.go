package gelatin

import (
	"fmt"
	"io"
)

type GelatinItemFilterName int

const (
	GelatinItemFilterFilters GelatinItemFilterName = iota
	GelatinItemFilterFiltersIsPlayed
	GelatinItemFilterFiltersIsFolder
	GelatinItemFilterFiltersIsNotFolder
	GelatinItemFilterParentId
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
	UnplayedItemCount     int32   // For series
	PlayedPercentage      float64 // For series
}

func (i *GelatinLibraryItemUserActivity) IsMatch(other *GelatinLibraryItemUserActivity) bool {
	if i.IsFavorite != other.IsFavorite {
		return false
	}
	if i.Played != other.Played {
		return false
	}
	if i.PlayCount != other.PlayCount {
		return false
	}
	if i.PlaybackPositionTicks != other.PlaybackPositionTicks {
		return false
	}
	if i.Rating != other.Rating {
		return false
	}
	if i.Rating != other.Rating {
		return false
	}

	return true
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
	Type         string // Movie, Series, Season, Episode, etc.
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
	// GetItems returns library items that match the given filters.
	//
	// If `recursive` is true, the search will recurse through library folders.
	//
	// Refer to Emby or Jellyfin docs for available item filters.
	GetItems(filters map[string]string, recursive bool) ([]GelatinLibraryItem, error)

	// GetItemsByUser returns library items for a _specific_ user (i.e., with user activity attached)
	//
	// Refer to Emby or Jellyfin docs for available item filters.
	GetItemsByUser(id string, filters map[string]string) ([]GelatinLibraryItem, error)

	// UpdateItem updates the given item
	//
	// Note that this does not modify item user data!
	UpdateItem(itemId string, item *GelatinLibraryItem) error

	// UpdateItemUserData updates the user data for the given item
	UpdateItemUserActivity(itemId string, userId string, old, new *GelatinLibraryItemUserActivity) error

	// GetItemFilterString returns the string representation of the given filter
	GetItemFilterString(filter GelatinItemFilterName) string
}

type GelatinPlaylistService interface {
}

type GelatinService interface {
	// ApiKey returns the current API key used by the client
	ApiKey() ApiKey

	// SetApiKey sets the API key
	SetApiKey(key ApiKey)

	System() GelatinSystemService
	User() GelatinUserService
	Library() GelatinLibraryService
	Playlist() GelatinPlaylistService
}

// Gets a user by name from the given service
func getUserByName(s GelatinService, username string) (*GelatinUser, error) {
	users, err := s.User().GetUsers(false)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.Name == username {
			return &user, nil
		}
	}

	return nil, fmt.Errorf("user %q not found", username)
}
