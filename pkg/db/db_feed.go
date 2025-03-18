package db

import (
	"context"
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type Feed struct {
	ID string `json:"id"`
	Name string `json:"name"`
	ProjectID string `json:"project_id"`
	UserID string `json:"user_id"`
	GenericDatabaseInfo
}

func (env *env) GetFeedsByProject(ctx context.Context, userID string, projectID string) ([]Feed, error) {

	env.logger.Info("Getting feeds by project")

	container, err := env.client.NewContainer("feeds")
	if err != nil {
		return []Feed{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	query := "SELECT * FROM feeds f WHERE f.user_id = @user_id AND f.project_id = @project_id"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
			{Name: "@project_id", Value: projectID},
		},
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	feeds := []Feed{}

	for pager.More() {
		response, err := pager.NextPage(ctx)
		if err != nil {
			return []Feed{}, err
		}

		for _, bytes := range response.Items {
			feed := Feed{}
			err := json.Unmarshal(bytes, &feed)
			if err != nil {
				return []Feed{}, err
			}
			feeds = append(feeds, feed)
		}
	}

	return feeds, nil

}

func (env *env) GetFeedByID(ctx context.Context, userID string, feedID string) (Feed, error) {

	env.logger.Info("Getting feed by ID")
	feed := Feed{}

	container, err := env.client.NewContainer("feeds")
	if err != nil {
		return feed, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	response, err := container.ReadItem(ctx, partitionKey, feedID, nil)
	if err != nil {
		return feed, err
	}

	err = json.Unmarshal(response.Value, &feed)
	if err != nil {
		return feed, err
	}

	return feed, nil

}

func (env *env) UpsertFeed(ctx context.Context, feed Feed) (Feed, error) {
	
	env.logger.Info("Upserting feed")

	container, err := env.client.NewContainer("feeds")

	partitionKey := azcosmos.NewPartitionKeyString(feed.UserID)

	marshalled, err := json.Marshal(feed)
	if err != nil {
		return feed, err
	}

	_, err = container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return feed, err
	}

	return feed, nil

}

func (env *env) RemoveFeed(ctx context.Context, userID string, feedID string) (interface{}, error) {

	env.logger.Info("Removing feed")

	container, err := env.client.NewContainer("feeds")

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	response, err := container.DeleteItem(ctx, partitionKey, feedID, nil)
	if err != nil {
		return nil, err
	}

	return response, nil

}