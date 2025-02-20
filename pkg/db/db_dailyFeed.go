package db

import (
	"context"
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type DailyFeed struct {
	ID string `json:"id"`
	FeedDate string `json:"feed_date"`
	FeedAmount float64 `json:"feed_amount"`
	AnimalID string `json:"animalid"`
	FeedID string `json:"feedid"`
	FeedPurchaseID string `json:"feedpurchaseid"`
	ProjectID string `json:"projectid"`
	UserID string `json:"userid"`
	GenericDatabaseInfo
}

func (env *env) GetDailyFeedsByProjectAndAnimal(ctx context.Context, userid string, projectid string, animalid string) ([]DailyFeed, error) {

	env.logger.Info("Getting daily feeds by project and animal")

	container, err := env.client.NewContainer("dailyfeeds")
	if err != nil {
		return []DailyFeed{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	query := "SELECT * FROM dailyfeeds df WHERE df.userid = @id AND df.projectid = @projectid AND df.animalid = @animalid"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@id", Value: userid},
			{Name: "@projectid", Value: projectid},
			{Name: "@animalid", Value: animalid},
		},
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	dailyFeeds := []DailyFeed{}

	for pager.More() {
		response, err := pager.NextPage(ctx)
		if err != nil {
			return []DailyFeed{}, err
		}

		for _, bytes := range response.Items {
			dailyFeed := DailyFeed{}
			err := json.Unmarshal(bytes, &dailyFeed)
			if err != nil {
				return []DailyFeed{}, err
			}
			dailyFeeds = append(dailyFeeds, dailyFeed)
		}
	}

	return dailyFeeds, nil

}

func (env *env) GetDailyFeedByID(ctx context.Context, userid string, dailyFeedID string) (DailyFeed, error) {

	env.logger.Info("Getting daily feed by ID")
	dailyFeed := DailyFeed{}

	container, err := env.client.NewContainer("dailyfeeds")
	if err != nil {
		return dailyFeed, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	response, err := container.ReadItem(ctx, partitionKey, dailyFeedID, nil)
	if err != nil {
		return dailyFeed, err
	}

	err = json.Unmarshal(response.Value, &dailyFeed)
	if err != nil {
		return dailyFeed, err
	}

	return dailyFeed, nil

}

func (env *env) UpsertDailyFeed(ctx context.Context, dailyFeed DailyFeed) (interface{}, error) {
	
	env.logger.Info("Upserting daily feed")

	container, err := env.client.NewContainer("dailyfeeds")

	partitionKey := azcosmos.NewPartitionKeyString(dailyFeed.UserID)

	marshalled, err := json.Marshal(dailyFeed)
	if err != nil {
		return nil, err
	}

	response, err := container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return nil, err
	}

	return response, nil

}

func (env *env) RemoveDailyFeed(ctx context.Context, userid string, dailyfeedid string) (interface{}, error) {

	env.logger.Info("Removing daily feed")

	container, err := env.client.NewContainer("dailyfeeds")

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	response, err := container.DeleteItem(ctx, partitionKey, dailyfeedid, nil)
	if err != nil {
		return nil, err
	}

	return response, nil

}