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
	AnimalID string `json:"animal_id"`
	FeedID string `json:"feed_id"`
	FeedPurchaseID string `json:"feed_purchase_id"`
	ProjectID string `json:"project_id"`
	UserID string `json:"user_id"`
	GenericDatabaseInfo
}

func (env *env) GetDailyFeedsByProjectAndAnimal(ctx context.Context, userID string, projectID string, animalID string) ([]DailyFeed, error) {

	env.logger.Info("Getting daily feeds by project and animal")

	container, err := env.client.NewContainer("dailyfeeds")
	if err != nil {
		return []DailyFeed{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	query := "SELECT * FROM dailyfeeds df WHERE df.user_id = @user_id AND df.project_id = @project_id AND df.animal_id = @animal_id"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
			{Name: "@project_id", Value: projectID},
			{Name: "@animal_id", Value: animalID},
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

func (env *env) GetDailyFeedByID(ctx context.Context, userID string, dailyFeedID string) (DailyFeed, error) {

	env.logger.Info("Getting daily feed by ID")
	dailyFeed := DailyFeed{}

	container, err := env.client.NewContainer("dailyfeeds")
	if err != nil {
		return dailyFeed, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

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

func (env *env) UpsertDailyFeed(ctx context.Context, dailyFeed DailyFeed) (DailyFeed, error) {
	
	env.logger.Info("Upserting daily feed")

	container, err := env.client.NewContainer("dailyfeeds")

	partitionKey := azcosmos.NewPartitionKeyString(dailyFeed.UserID)

	marshalled, err := json.Marshal(dailyFeed)
	if err != nil {
		return dailyFeed, err
	}

	_, err = container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return dailyFeed, err
	}

	return dailyFeed, nil

}

func (env *env) RemoveDailyFeed(ctx context.Context, userID string, dailyFeedID string) (interface{}, error) {

	env.logger.Info("Removing daily feed")

	container, err := env.client.NewContainer("dailyfeeds")

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	response, err := container.DeleteItem(ctx, partitionKey, dailyFeedID, nil)
	if err != nil {
		return nil, err
	}

	return response, nil

}