package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type Expense struct {
	ID        string  `json:"id"`
	Date      string  `json:"date"`
	Items     string  `json:"items"`
	Quantity  float64 `json:"quantity"`
	Cost      float64 `json:"cost"`
	ProjectID string  `json:"project_id"`
	UserID    string  `json:"user_id"`
	GenericDatabaseInfo
}

func (e Expense) GetID() string {
	return e.ID
}

func (env *env) GetExpensesByProject(ctx context.Context, userID string, projectID string, paginationOptions PaginationOptions) ([]Expense, error) {

	env.logger.Info("Getting expenses by project")

	container, err := env.client.NewContainer("expenses")
	if err != nil {
		return []Expense{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	sortOrder := "ASC"
	if paginationOptions.SortByNewest {
		sortOrder = "DESC"
	}

	query := fmt.Sprintf("SELECT * FROM expenses e WHERE e.user_id = @user_id AND e.project_id = @project_id ORDER BY e.created %s", sortOrder)

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
			{Name: "@project_id", Value: projectID},
		},
		PageSizeHint: int32(paginationOptions.PerPage),
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	expenses := []Expense{}
	currentPage := 0

	for pager.More() {

		if currentPage == paginationOptions.Page {
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

			return expenses, nil

		} else {
			_, err := pager.NextPage(ctx)
			if err != nil {
				return []Expense{}, err
			}
			currentPage++
		}

	}

	return expenses, nil

}

func (env *env) GetProjectDependentExpenses(ctx context.Context, userID string, projectID string) ([]Identifiable, error) {

	env.logger.Info("Getting project dependent expenses")

	container, err := env.client.NewContainer("expenses")
	if err != nil {
		return []Identifiable{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	query := "SELECT * FROM expenses e WHERE e.user_id = @user_id AND e.project_id = @project_id"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
			{Name: "@project_id", Value: projectID},
		},
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	expenses := []Expense{}

	for pager.More() {

		response, err := pager.NextPage(ctx)
		if err != nil {
			return []Identifiable{}, err
		}

		for _, bytes := range response.Items {
			expense := Expense{}
			err := json.Unmarshal(bytes, &expense)
			if err != nil {
				return []Identifiable{}, err
			}
			expenses = append(expenses, expense)
		}

	}

	identifiables := []Identifiable{}

	for _, e := range expenses {
		identifiables = append(identifiables, e)
	}

	return identifiables, nil

}

func (env *env) GetExpenseByID(ctx context.Context, userID string, expenseID string) (Expense, error) {

	env.logger.Info("Getting expense by ID")
	expense := Expense{}

	container, err := env.client.NewContainer("expenses")
	if err != nil {
		return expense, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	response, err := container.ReadItem(ctx, partitionKey, expenseID, nil)
	if err != nil {
		return expense, err
	}

	err = json.Unmarshal(response.Value, &expense)
	if err != nil {
		return expense, err
	}

	return expense, nil

}

func (env *env) UpsertExpense(ctx context.Context, expense Expense) (Expense, error) {

	env.logger.Info("Upserting expense")

	container, err := env.client.NewContainer("expenses")
	if err != nil {
		return expense, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(expense.UserID)

	marshalled, err := json.Marshal(expense)
	if err != nil {
		return expense, err
	}

	_, err = container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return expense, err
	}

	return expense, nil

}

func (env *env) RemoveExpense(ctx context.Context, userID string, expenseID string) (interface{}, error) {

	env.logger.Info("Removing expense")

	container, err := env.client.NewContainer("expenses")
	if err != nil {
		return nil, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	response, err := container.DeleteItem(ctx, partitionKey, expenseID, nil)
	if err != nil {
		return nil, err
	}

	return response, nil

}
