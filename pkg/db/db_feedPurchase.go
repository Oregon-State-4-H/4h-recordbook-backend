package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type FeedPurchase struct {
	ID              string  `json:"id"`
	DatePurchased   string  `json:"date_purchased"`
	AmountPurchased float64 `json:"amount_purchased"`
	TotalCost       float64 `json:"total_cost"`
	FeedID          string  `json:"feed_id"`
	ProjectID       string  `json:"project_id"`
	UserID          string  `json:"user_id"`
	GenericDatabaseInfo
}

func (fp FeedPurchase) GetID() string {
	return fp.ID
}

func (env *env) GetFeedPurchasesByProject(ctx context.Context, userID string, projectID string, paginationOptions PaginationOptions) ([]FeedPurchase, error) {

	env.logger.Info("Getting feed purchases by project")

	container, err := env.client.NewContainer("feedpurchases")
	if err != nil {
		return []FeedPurchase{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	sortOrder := "ASC"
	if paginationOptions.SortByNewest {
		sortOrder = "DESC"
	}

	query := fmt.Sprintf("SELECT * FROM feedpurchases fp WHERE fp.user_id = @user_id AND fp.project_id = @project_id ORDER BY fp.created %s", sortOrder)

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
			{Name: "@project_id", Value: projectID},
		},
		PageSizeHint: int32(paginationOptions.PerPage),
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	feedPurchases := []FeedPurchase{}
	currentPage := 0

	for pager.More() {

		if currentPage == paginationOptions.Page {
			response, err := pager.NextPage(ctx)
			if err != nil {
				return []FeedPurchase{}, err
			}

			for _, bytes := range response.Items {
				feedPurchase := FeedPurchase{}
				err := json.Unmarshal(bytes, &feedPurchase)
				if err != nil {
					return []FeedPurchase{}, err
				}
				feedPurchases = append(feedPurchases, feedPurchase)
			}

			return feedPurchases, nil

		} else {
			_, err := pager.NextPage(ctx)
			if err != nil {
				return []FeedPurchase{}, err
			}
			currentPage++
		}

	}

	return feedPurchases, nil

}

func (env *env) GetFeedDependentFeedPurchases(ctx context.Context, userID string, feedID string) ([]Identifiable, error) {

	env.logger.Info("Getting feed dependent feed purchases")

	container, err := env.client.NewContainer("feedpurchases")
	if err != nil {
		return []Identifiable{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	query := "SELECT * FROM feedpurchases fp WHERE fp.user_id = @user_id AND fp.feed_id = @feed_id"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
			{Name: "@feed_id", Value: feedID},
		},
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	feedPurchases := []FeedPurchase{}

	for pager.More() {

		response, err := pager.NextPage(ctx)
		if err != nil {
			return []Identifiable{}, err
		}

		for _, bytes := range response.Items {
			feedPurchase := FeedPurchase{}
			err := json.Unmarshal(bytes, &feedPurchase)
			if err != nil {
				return []Identifiable{}, err
			}
			feedPurchases = append(feedPurchases, feedPurchase)
		}

	}

	identifiables := []Identifiable{}

	for _, fp := range feedPurchases {
		identifiables = append(identifiables, fp)
	}

	return identifiables, nil

}

func (env *env) GetFeedPurchaseByID(ctx context.Context, userID string, feedPurchaseID string) (FeedPurchase, error) {

	env.logger.Info("Getting feed purchase by ID")
	feedPurchase := FeedPurchase{}

	container, err := env.client.NewContainer("feedpurchases")
	if err != nil {
		return feedPurchase, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	response, err := container.ReadItem(ctx, partitionKey, feedPurchaseID, nil)
	if err != nil {
		return feedPurchase, err
	}

	err = json.Unmarshal(response.Value, &feedPurchase)
	if err != nil {
		return feedPurchase, err
	}

	return feedPurchase, nil

}

func (env *env) UpsertFeedPurchase(ctx context.Context, feedPurchase FeedPurchase) (FeedPurchase, error) {

	env.logger.Info("Upserting feed purchase")

	container, err := env.client.NewContainer("feedpurchases")
	if err != nil {
		return feedPurchase, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(feedPurchase.UserID)

	marshalled, err := json.Marshal(feedPurchase)
	if err != nil {
		return feedPurchase, err
	}

	_, err = container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return feedPurchase, err
	}

	return feedPurchase, nil

}

func (env *env) RemoveFeedPurchase(ctx context.Context, userID string, feedPurchaseID string) (interface{}, error) {

	env.logger.Info("Removing feed purchase")

	container, err := env.client.NewContainer("feedpurchases")
	if err != nil {
		return nil, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	response, err := container.DeleteItem(ctx, partitionKey, feedPurchaseID, nil)
	if err != nil {
		return nil, err
	}

	return response, nil

}
