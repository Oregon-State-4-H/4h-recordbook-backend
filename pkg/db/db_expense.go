package db

import (
	"context"
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type Expense struct {
	ID string `json:"id"`
	Date string `json:"date"`
	Items string `json:"items"`
	Quantity float64 `json:"quantity"`
	Cost float64 `json:"cost"`
	ProjectID string `json:"projectid"`
	UserID string `json:"userid"`
	GenericDatabaseInfo
}

func (env *env) GetExpensesByProject(ctx context.Context, userid string, projectid string) ([]Expense, error) {

	env.logger.Info("Getting expenses by project")

	container, err := env.client.NewContainer("expenses")
	if err != nil {
		return []Expense{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	query := "SELECT * FROM expenses e WHERE e.userid = @id AND e.projectid = @projectid"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@id", Value: userid},
			{Name: "@projectid", Value: projectid},
		},
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	expenses := []Expense{}

	for pager.More() {
		response, err := pager.NextPage(ctx)
		if err != nil {
			return []Expense{}, err
		}

		for _, bytes := range response.Items {
			expense := Expense{}
			err := json.Unmarshal(bytes, &expense)
			if err != nil {
				return []Expense{}, err
			}
			expenses = append(expenses, expense)
		}
	}

	return expenses, nil

}

func (env *env) GetExpenseByID(ctx context.Context, userid string, expenseid string) (Expense, error) {

	env.logger.Info("Getting expense by ID")
	expense := Expense{}

	container, err := env.client.NewContainer("expenses")
	if err != nil {
		return expense, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	response, err := container.ReadItem(ctx, partitionKey, expenseid, nil)
	if err != nil {
		return expense, err
	}

	err = json.Unmarshal(response.Value, &expense)
	if err != nil {
		return expense, err
	}

	return expense, nil

}

func (env *env) UpsertExpense(ctx context.Context, expense Expense) (interface{}, error) {
	
	env.logger.Info("Upserting expense")

	container, err := env.client.NewContainer("expenses")

	partitionKey := azcosmos.NewPartitionKeyString(expense.UserID)

	marshalled, err := json.Marshal(expense)
	if err != nil {
		return nil, err
	}

	response, err := container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return nil, err
	}

	return response, nil

}