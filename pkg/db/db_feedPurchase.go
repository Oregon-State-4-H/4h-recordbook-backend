package db

import (
	"context"
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type FeedPurchase struct {
	ID string `json:"id"`
	DatePurchased string `json:"date_purchased"`
	AmountPurchased float64 `json:"amount_purchased"`
	TotalCost float64 `json:"total_cost"`
	FeedID string `json:"feedid"`
	ProjectID string `json:"projectid"`
	UserID string `json:"userid"`
	GenericDatabaseInfo
}

func (env *env) GetFeedPurchasesByProject(ctx context.Context, userid string, projectid string) ([]FeedPurchase, error) {

	env.logger.Info("Getting feed purchases by project")

	container, err := env.client.NewContainer("feedpurchases")
	if err != nil {
		return []FeedPurchase{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	query := "SELECT * FROM feedpurchases fp WHERE fp.userid = @id AND fp.projectid = @projectid"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@id", Value: userid},
			{Name: "@projectid", Value: projectid},
		},
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	feedPurchases := []FeedPurchase{}

	for pager.More() {
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
	}

	return feedPurchases, nil

}

func (env *env) GetFeedPurchaseByID(ctx context.Context, userid string, feedPurchaseID string) (FeedPurchase, error) {

	env.logger.Info("Getting feed purchase by ID")
	feedPurchase := FeedPurchase{}

	container, err := env.client.NewContainer("feedpurchases")
	if err != nil {
		return feedPurchase, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

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

func (env *env) UpsertFeedPurchase(ctx context.Context, feedPurchase FeedPurchase) (interface{}, error) {
	
	env.logger.Info("Upserting feed purchase")

	container, err := env.client.NewContainer("feedpurchases")

	partitionKey := azcosmos.NewPartitionKeyString(feedPurchase.UserID)

	marshalled, err := json.Marshal(feedPurchase)
	if err != nil {
		return nil, err
	}

	response, err := container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return nil, err
	}

	return response, nil

}