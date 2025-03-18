package db

import (
	"context"
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type Supply struct {
	ID string `json:"id"`
	Description string `json:"description"`
	StartValue float64 `json:"start_value"`
	EndValue float64 `json:"end_value"`
	ProjectID string `json:"project_id"`
	UserID string `json:"user_id"`
	GenericDatabaseInfo
}

func (env *env) GetSuppliesByProject(ctx context.Context, userID string, projectID string) ([]Supply, error) {

	env.logger.Info("Getting supplies by project")

	container, err := env.client.NewContainer("supplies")
	if err != nil {
		return []Supply{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	query := "SELECT * FROM supplies s WHERE s.user_id = @user_id AND s.project_id = @project_id"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
			{Name: "@project_id", Value: projectID},
		},
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	supplies := []Supply{}

	for pager.More() {
		response, err := pager.NextPage(ctx)
		if err != nil {
			return []Supply{}, err
		}

		for _, bytes := range response.Items {
			supply := Supply{}
			err := json.Unmarshal(bytes, &supply)
			if err != nil {
				return []Supply{}, err
			}
			supplies = append(supplies, supply)
		}
	}

	return supplies, nil

}

func (env *env) GetSupplyByID(ctx context.Context, userID string, supplyID string) (Supply, error) {

	env.logger.Info("Getting supply by ID")
	supply := Supply{}

	container, err := env.client.NewContainer("supplies")
	if err != nil {
		return supply, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	response, err := container.ReadItem(ctx, partitionKey, supplyID, nil)
	if err != nil {
		return supply, err
	}

	err = json.Unmarshal(response.Value, &supply)
	if err != nil {
		return supply, err
	}

	return supply, nil

}

func (env *env) UpsertSupply(ctx context.Context, supply Supply) (Supply, error) {
	
	env.logger.Info("Upserting supply")

	container, err := env.client.NewContainer("supplies")

	partitionKey := azcosmos.NewPartitionKeyString(supply.UserID)

	marshalled, err := json.Marshal(supply)
	if err != nil {
		return supply, err
	}

	_, err = container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return supply, err
	}

	return supply, nil

}

func (env *env) RemoveSupply(ctx context.Context, userID string, supplyID string) (interface{}, error) {

	env.logger.Info("Removing supply")

	container, err := env.client.NewContainer("supplies")

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	response, err := container.DeleteItem(ctx, partitionKey, supplyID, nil)
	if err != nil {
		return nil, err
	}

	return response, nil

}