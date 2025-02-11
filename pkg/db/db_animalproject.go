package db

import (
	"context"
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type Animal struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Species string `json:"species"`
	BirthDate string `json:"birth_date"`
	PurchaseDate string `json:"purchase_date"`
	SireBreed string `json:"sire_breed"`
	DamBreed string `json:"dam_breed"`
	BeginningWeight float64 `json:"beginning_weight"`
	EndWeight float64 `json:"end_weight"`
	EndDate string `json:"end_date"`
	AnimalCost string `json:"animal_cost"`
	SalePrice string `json:"sale_price"`
	YieldGrade string `json:"yield_grade"`
	QualityGrade string `json:"quality_grade"`
	UserID string `json:"userid"`
	GenericDatabaseInfo
}

type AnimalProjectRelation struct {
	ID string `json:"id"`
	AnimalID string `json:"animalid"`
	ProjectID string `json:"projectid"`
	GenericDatabaseInfo
}

type Feed struct {
	ID string `json:"id"`
	Name string `json:"name"`
	ProjectID string `json:"projectid"`
	UserID string `json:"userid"`
	GenericDatabaseInfo
}

type FeedPurchase struct {
	ID string `json:"id"`
	DatePurchased string `json:"date_purchased"`
	AmountPurchased float64 `json:"amount_purchased"`
	TotalCost float64 `json:"total_cost"`
	FeedID string `json:"feedid"`
	ProjectID string `json:"projectid`
	UserID string `json:"userid"`
	GenericDatabaseInfo
}

type DailyFeed struct {
	ID string `json:"id"`
	FeedDate string `json:"feed_date"`
	FeedAmount float64 `json:"feed_amount"`
	AnimalID string `json:"animalid"`
	FeedID string `json:"feedid"`
	FeedPurchaseID string `json:"feedpurchaseid"`
	ProjectID string `json:"projectid`
	UserID string `json:"userid"`
	GenericDatabaseInfo
}

type Expenses struct {
	ID string `json:"id"`
	Date string `json:"date"`
	Items string `json:"items"`
	Quantity float64 `json:"quantity"`
	Cost float64 `json:"cost"`
	ProjectID string `json:"projectid`
	UserID string `json:"userid"`
	GenericDatabaseInfo
}

type Supplies struct {
	ID string `json:"id"`
	Description string `json:"description"`
	StartValue float64 `json:"start_value"`
	EndValue float64 `json:"end_value"`
	ProjectID string `json:"projectid`
	UserID string `json:"userid"`
	GenericDatabaseInfo
}

/*******************************
* ANIMALS
********************************/

/*******************************
* FEED
********************************/
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

/*******************************
* DAILY FEED
********************************/

/*******************************
* EXPENSES
********************************/

/*******************************
* SUPPLY INVENTORY
********************************/