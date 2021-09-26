package gelatin

import (
	"fmt"
	"log"
	"strings"

	"github.com/google/go-cmp/cmp"
)

type GelatinClientOpts struct {
	Interactive bool
}

type GelatinClient struct {
	// Export data from this service
	from GelatinService

	// Import data into this service
	into GelatinService

	opts GelatinClientOpts
}

func NewGelatinClient(from GelatinService, into GelatinService, opts *GelatinClientOpts) *GelatinClient {
	if from == nil || into == nil {
		return nil
	}

	if opts != nil {
		return &GelatinClient{from: from, into: into, opts: *opts}
	} else {
		return &GelatinClient{from: from, into: into}
	}
}

func (c *GelatinClient) MigrateUsers(passwords map[string]string) error {
	fromUsers, err := c.from.User().GetUsers(false)
	if err != nil {
		return err
	}

	intoUsers, err := c.into.User().GetUsers(false)
	if err != nil {
		return err
	}

	var createUsers []*GelatinUser
	var deleteUsers []*GelatinUser

	// Find users we need to create
	for i, fromUser := range fromUsers {
		exists := false
		for _, intoUser := range intoUsers {
			if fromUser.Name == intoUser.Name {
				exists = true
				break
			}
		}

		if !exists {
			createUsers = append(createUsers, &fromUsers[i])
		}
	}

	// Find users we need to delete
	for i, intoUser := range intoUsers {
		exists := false
		for _, fromUser := range fromUsers {
			if intoUser.Name == fromUser.Name {
				exists = true
				break
			}
		}

		if !exists {
			deleteUsers = append(deleteUsers, &intoUsers[i])
		}
	}

	for i, user := range createUsers {
		newUser, err := c.into.User().CreateUser(user.Name)
		if err != nil {
			return err
		}

		log.Printf("created %s: %s", newUser.Name, newUser.Id)

		createUsers[i] = newUser
	}

	// TODO: Migrate this to iterate over deleteUsers
	for i, user := range createUsers {
		err := c.into.User().DeleteUser(user.Id)
		if err != nil {
			return err
		}

		log.Printf("Deleted: %s", user.Name)

		createUsers[i] = nil
	}

	return nil
}

// DiffUsers returns a diff of the list of users
//
// If full is true, the diff will include the full user struct.
func (c *GelatinClient) DiffUsers(full bool) (string, error) {
	fromUsers, err := c.from.User().GetUsers(false)
	if err != nil {
		return "", err
	}

	intoUsers, err := c.into.User().GetUsers(false)
	if err != nil {
		return "", err
	}

	if full {
		return cmp.Diff(fromUsers, intoUsers), nil
	}

	var fromUsernames []string
	var intoUsernames []string

	for _, user := range fromUsers {
		fromUsernames = append(fromUsernames, user.Name)
	}

	for _, user := range intoUsers {
		intoUsernames = append(intoUsernames, user.Name)
	}

	return cmp.Diff(fromUsernames, intoUsernames), nil
}

// getProviderIds returns a list of provider IDs for an item
func getProviderIds(item *GelatinLibraryItem) []string {
	var providerIds []string
	if item.ImdbId != "" {
		providerIds = append(providerIds, item.ImdbId)
	}
	if item.TmdbId != "" {
		providerIds = append(providerIds, item.TmdbId)
	}
	if item.TvdbId != "" {
		providerIds = append(providerIds, item.TvdbId)
	}

	// As a fallback, use the name of the item as a provider
	if len(providerIds) == 0 {
		providerIds = append(providerIds, item.Name)
	}

	return providerIds
}

// handleSeries recursively updates the playedItemData map for a given series.
//
// It works like this:
//
// a. Get the series' provider ID
// b. If the series is fully played, store an entry for the series itself and skip the remaining steps
// c. Run a second query for the series' seasons (children)
// d. For each season, check if it is fully played; if it is, store an entry and return
// e. If the current season is not fully played, walk through each episode of the season and store an individual entry for the episode
func handleSeries(svc GelatinLibraryService, item *GelatinLibraryItem, userId string, seriesProviderIds []string, playedItemData map[string]*GelatinLibraryItem) error {
	switch item.Type {
	case "Series":
		if item.UserData.Played && item.UserData.PlayedPercentage == 100 {
			// The series is fully played, so just add an entry for it
			for _, id := range getProviderIds(item) {
				playedItemData[id] = item
			}
		} else {
			seriesProviderIds := getProviderIds(item)

			children, err := svc.GetItemsByUser(userId, map[string]string{
				svc.GetItemFilterString(GelatinItemFilterParentId): item.Id,
			})
			if err != nil {
				return err
			}

			// These could either be episodes or seasons
			for i := range children {
				handleSeries(svc, &children[i], userId, seriesProviderIds, playedItemData)
			}
		}
	case "Season":
		if item.UserData.Played && item.UserData.PlayedPercentage == 100 {
			// The season is fully played, so just add an entry for it (for each provider ID)
			for _, id := range seriesProviderIds {
				seasonKey := fmt.Sprintf("%s-%d", id, item.IndexNumber)
				playedItemData[seasonKey] = item
			}
		} else {
			episodes, err := svc.GetItemsByUser(userId, map[string]string{
				svc.GetItemFilterString(GelatinItemFilterParentId): item.Id,
			})
			if err != nil {
				return err
			}

			for i := range episodes {
				handleSeries(svc, &episodes[i], userId, seriesProviderIds, playedItemData)
			}
		}
	case "Episode":
		for _, id := range seriesProviderIds {
			episodeKey := fmt.Sprintf("%s-%d-%d", id, item.ParentIndexNumber, item.IndexNumber)
			playedItemData[episodeKey] = item
		}
	}

	return nil
}

func promptUserYesNo(message string) bool {
	fmt.Printf("%s (y/n) [y]: ", message)
	var in string
	fmt.Scanf("%s", &in)
	return strings.ToLower(strings.TrimSpace(in)) != "n"
}

func promptUserForWatchItem(old, new *GelatinLibraryItem) bool {
	var message string

	oldPlayed, newPlayed := old.UserData.Played, new.UserData.Played
	oldFavorite, newFavorite := old.UserData.IsFavorite, new.UserData.IsFavorite
	oldTicks, newTicks := old.UserData.PlaybackPositionTicks, new.UserData.PlaybackPositionTicks

	var played, favorite, ticks string

	if oldPlayed != newPlayed {
		played = fmt.Sprintf("%v->%v", oldPlayed, newPlayed)
	} else {
		played = fmt.Sprintf("%v", oldPlayed)
	}

	if oldFavorite != newFavorite {
		favorite = fmt.Sprintf("%v->%v", oldFavorite, newFavorite)
	} else {
		favorite = fmt.Sprintf("%v", oldFavorite)
	}

	if oldTicks != newTicks {
		ticks = fmt.Sprintf("%v->%v", oldTicks, newTicks)
	} else {
		ticks = fmt.Sprintf("%v", oldTicks)
	}

	switch old.Type {
	case "Movie":
		message = fmt.Sprintf("Movie: %q, Played: %s, Favorite: %s, Ticks: %s", old.Name, played, favorite, ticks)
	case "Series":
		message = fmt.Sprintf("Series: %q, Played: %s, Favorite: %s", old.Name, played, favorite)
	case "Season":
		message = fmt.Sprintf("Season: %q (%q), Played: %s, Favorite: %s", old.Name, old.SeriesName, played, favorite)
	case "Episode":
		message = fmt.Sprintf("Episode: %q (%q), Played: %s, Favorite: %s, Ticks: %s", old.Name, old.SeriesName, played, favorite, ticks)
	default:
		panic(fmt.Sprintf("unexpected item type: %s for %v", old.Type, old))
	}

	return promptUserYesNo(message)
}

// MigrateUserWatchHistory migrates a user's watch history from one service to another.
//
// If the user does not exist in either service, this method returns an error.
//
// Migration works like this:
//
// 1. Fetch all movies and series from the "from" service.
// 2. For each movie and series/episode, store an entry containing the user activity using the provider ID (IMDb, TMDB, TVDB)
// 3. Fetch all items from the into service and compare the user activity state with that of the from service
// 4. If there is a difference, update the into service with the latest state
func (c *GelatinClient) MigrateUserWatchHistory(username string) error {
	fromUser, err := getUserByName(c.from, username)
	if err != nil {
		return err
	}

	intoUser, err := getUserByName(c.into, username)
	if err != nil {
		return err
	}

	// Get all items for the user in the from service
	fromLibraryItems, err := c.from.Library().GetItemsByUser(fromUser.Id, nil)
	if err != nil {
		return err
	}

	// Build a map of user data for the items played in the from service
	playedItemData := make(map[string]*GelatinLibraryItem)
	for i, item := range fromLibraryItems {
		switch item.Type {
		case "Movie":
			for _, id := range getProviderIds(&item) {
				playedItemData[id] = &fromLibraryItems[i]
			}
		case "Series":
			// Recursively handle this series
			err := handleSeries(c.from.Library(), &fromLibraryItems[i], fromUser.Id, nil, playedItemData)
			if err != nil {
				return err
			}
		}
	}

	// Get all library items tracked by the into service
	intoLibraryItems, err := c.into.Library().GetItemsByUser(intoUser.Id, nil)
	if err != nil {
		return err
	}

	// Do a first pass to build a mapping from series ID to known provider IDs
	seriesIdToProviderIds := make(map[string][]string)
	for _, item := range intoLibraryItems {
		if item.Type == "Series" {
			seriesIdToProviderIds[item.Id] = getProviderIds(&item)
		}
	}

	// Finally, run through the list of items once more and update user watch state if it differs
	for _, item := range intoLibraryItems {
		var fromItem *GelatinLibraryItem

		switch item.Type {
		case "Movie", "Series":
			keys := []string{item.ImdbId, item.TmdbId, item.TvdbId}
			for _, key := range keys {
				if data, ok := playedItemData[key]; ok {
					fromItem = data
					break
				}
			}
		case "Season":
			providerIds := seriesIdToProviderIds[item.SeriesId]
			for _, id := range providerIds {
				key := fmt.Sprintf("%s-%d", id, item.IndexNumber)
				if data, ok := playedItemData[key]; ok {
					fromItem = data
					break
				}
			}
		case "Episode":
			providerIds := seriesIdToProviderIds[item.SeriesId]
			for _, id := range providerIds {
				key := fmt.Sprintf("%s-%d-%d", id, item.ParentIndexNumber, item.IndexNumber)
				if data, ok := playedItemData[key]; ok {
					fromItem = data
					break
				}
			}
		}

		if fromItem == nil {
			// This item is not present in the from service, so we can skip it
			continue
		}

		needUpdate := !item.UserData.IsMatch(fromItem.UserData)

		if needUpdate {
			if c.opts.Interactive {
				if !promptUserForWatchItem(&item, fromItem) {
					// User skipped this item
					continue
				}
			}

			err := c.into.Library().UpdateItemUserActivity(item.Id, intoUser.Id, item.UserData, fromItem.UserData)
			if err != nil {
				return fmt.Errorf("failed to set user data for item %q: %v", item.Name, err)
			}
		}
	}

	return nil
}
