package db

import (
	"context"
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type Bookmark struct {
	ID string `json:"id"`
	Link string	`json:"link"`
	Label string `json:"label"`
	UserID string `json:"userid"`
	GenericDatabaseInfo
}

func (env *env) GetBookmarkByLink(ctx context.Context, userID string, link string) (Bookmark, error) {

	env.logger.Info("Getting bookmark by link")

	container, err := env.client.NewContainer("bookmarks")
	if err != nil {
		return Bookmark{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	query := "SELECT * FROM bookmarks b WHERE b.userid = @user_id AND b.link = @link"

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

	if len(bookmarks) == 0{
		return Bookmark{}, err
	}

	return bookmarks[0], nil

}

func (env *env) GetBookmarks(ctx context.Context, userID string) ([]Bookmark, error) {

	env.logger.Info("Getting bookmarks")

	container, err := env.client.NewContainer("bookmarks")
	if err != nil {
		return []Bookmark{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	query := "SELECT * FROM bookmarks b WHERE b.userid = @user_id"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
		},
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	bookmarks := []Bookmark{}

	for pager.More() {
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

	}

	return bookmarks, nil

} 

func (env *env) AddBookmark(ctx context.Context, bookmark Bookmark) (interface{}, error) {
	
	env.logger.Info("Adding bookmark")

	container, err := env.client.NewContainer("bookmarks")

	partitionKey := azcosmos.NewPartitionKeyString(bookmark.UserID)

	marshalled, err := json.Marshal(bookmark)
	if err != nil {
		return nil, err
	}

	response, err := container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return nil, err
	}

	return response, nil

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