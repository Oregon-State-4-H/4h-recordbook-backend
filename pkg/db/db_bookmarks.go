package db

import (
	"context"
	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type Bookmark struct {
	ID     string `json:"id"`
	Link   string `json:"link"`
	Label  string `json:"label"`
	UserID string `json:"user_id"`
	GenericDatabaseInfo
}

func (env *env) GetBookmarkByLink(ctx context.Context, userID string, link string) (Bookmark, error) {

	env.logger.Info("Getting bookmark by link")

	container, err := env.client.NewContainer("bookmarks")
	if err != nil {
		return Bookmark{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	query := "SELECT * FROM bookmarks b WHERE b.user_id = @user_id AND b.link = @link"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
			{Name: "@link", Value: link},
		},
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	bookmarks := []Bookmark{}

	for pager.More() {
		response, err := pager.NextPage(ctx)
		if err != nil {
			return Bookmark{}, err
		}

		for _, bytes := range response.Items {
			bookmark := Bookmark{}
			err := json.Unmarshal(bytes, &bookmark)
			if err != nil {
				return Bookmark{}, err
			}
			bookmarks = append(bookmarks, bookmark)
		}
	}

	if len(bookmarks) == 0 {
		return Bookmark{}, err
	}

	return bookmarks[0], nil

}

func (env *env) GetBookmarks(ctx context.Context, userID string, page int, perPage int) ([]Bookmark, error) {

	env.logger.Info("Getting bookmarks")

	container, err := env.client.NewContainer("bookmarks")
	if err != nil {
		return []Bookmark{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	query := "SELECT * FROM bookmarks b WHERE b.user_id = @user_id ORDER BY b.created ASC"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
		},
		PageSizeHint: int32(perPage),
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	bookmarks := []Bookmark{}
	currentPage := 0

	for pager.More() {

		if currentPage == page {
			response, err := pager.NextPage(ctx)
			if err != nil {
				return []Bookmark{}, err
			}

			for _, bytes := range response.Items {
				bookmark := Bookmark{}
				err := json.Unmarshal(bytes, &bookmark)
				if err != nil {
					return []Bookmark{}, err
				}
				bookmarks = append(bookmarks, bookmark)
			}

			return bookmarks, nil

		} else {
			_, err := pager.NextPage(ctx)
			if err != nil {
				return []Bookmark{}, err
			}
			currentPage++
		}

	}

	return bookmarks, nil

}

func (env *env) AddBookmark(ctx context.Context, bookmark Bookmark) (Bookmark, error) {

	env.logger.Info("Adding bookmark")

	container, err := env.client.NewContainer("bookmarks")

	partitionKey := azcosmos.NewPartitionKeyString(bookmark.UserID)

	marshalled, err := json.Marshal(bookmark)
	if err != nil {
		return bookmark, err
	}

	_, err = container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return bookmark, err
	}

	return bookmark, nil

}

func (env *env) RemoveBookmark(ctx context.Context, userID string, bookmarkID string) (interface{}, error) {

	env.logger.Info("Removing bookmark")

	container, err := env.client.NewContainer("bookmarks")

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	response, err := container.DeleteItem(ctx, partitionKey, bookmarkID, nil)
	if err != nil {
		return nil, err
	}

	return response, nil

}
