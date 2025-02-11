package db

import (
	"context"
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type Feed struct {
	ID string `json:"id"`
	Name string `json:"name"`
	ProjectID string `json:"projectid"`
	UserID string `json:"userid"`
	GenericDatabaseInfo
}

func (env *env) GetFeedsByProject(ctx context.Context, userid string, projectid string) ([]Feed, error) {

	env.logger.Info("Getting feeds by project")

	container, err := env.client.NewContainer("feeds")
	if err != nil {
		return []Feed{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	query := "SELECT * FROM feeds f WHERE f.userid = @id AND f.projectid = @projectid"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@id", Value: userid},
			{Name: "@projectid", Value: projectid},
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

func (env *env) GetFeedByID(ctx context.Context, userid string, feedid string) (Feed, error) {

	env.logger.Info("Getting feed by ID")
	feed := Feed{}

	container, err := env.client.NewContainer("feeds")
	if err != nil {
		return feed, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	response, err := container.ReadItem(ctx, partitionKey, feedid, nil)
	if err != nil {
		return feed, err
	}

	err = json.Unmarshal(response.Value, &feed)
	if err != nil {
		return feed, err
	}

	return feed, nil

}

func (env *env) UpsertFeed(ctx context.Context, feed Feed) (interface{}, error) {
	
	env.logger.Info("Upserting feed")

	container, err := env.client.NewContainer("feeds")

	partitionKey := azcosmos.NewPartitionKeyString(feed.UserID)

	marshalled, err := json.Marshal(feed)
	if err != nil {
		return nil, err
	}

	response, err := container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return nil, err
	}

	return response, nil

}