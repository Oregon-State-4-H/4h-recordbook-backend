package db

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type Project struct {
	ID          string `json:"id"`
	Year        string `json:"year"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	UserID      string `json:"user_id"`
	GenericDatabaseInfo
}

func (env *env) GetProjectByID(ctx context.Context, userID string, projectID string) (Project, error) {

	env.logger.Info("Getting project by ID")
	project := Project{}

	container, err := env.client.NewContainer("projects")
	if err != nil {
		return project, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	response, err := container.ReadItem(ctx, partitionKey, projectID, nil)
	if err != nil {
		return project, err
	}

	err = json.Unmarshal(response.Value, &project)
	if err != nil {
		return project, err
	}

	return project, nil

}

func (env *env) GetCurrentProjects(ctx context.Context, userID string, paginationOptions PaginationOptions) ([]Project, error) {

	env.logger.Info("Getting current projects")

	container, err := env.client.NewContainer("projects")
	if err != nil {
		return []Project{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	now := time.Now()
	year := strconv.Itoa(now.Year())

	sortOrder := "ASC"
	if paginationOptions.SortByNewest {
		sortOrder = "DESC"
	}

	query := fmt.Sprintf("SELECT * FROM projects p WHERE p.user_id = @user_id AND p.year = @year ORDER BY p.created %s", sortOrder)

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
			{Name: "@year", Value: year},
		},
		PageSizeHint: int32(paginationOptions.PerPage),
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	projects := []Project{}
	currentPage := 0

	for pager.More() {

		if currentPage == paginationOptions.Page {
			response, err := pager.NextPage(ctx)
			if err != nil {
				return []Project{}, err
			}

			for _, bytes := range response.Items {
				project := Project{}
				err := json.Unmarshal(bytes, &project)
				if err != nil {
					return []Project{}, err
				}
				projects = append(projects, project)
			}

			return projects, nil

		} else {
			_, err := pager.NextPage(ctx)
			if err != nil {
				return []Project{}, err
			}
			currentPage++
		}

	}

	return projects, nil

}

func (env *env) GetProjectsByUser(ctx context.Context, userID string, paginationOptions PaginationOptions) ([]Project, error) {

	env.logger.Info("Getting projects")

	container, err := env.client.NewContainer("projects")
	if err != nil {
		return []Project{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	sortOrder := "ASC"
	if paginationOptions.SortByNewest {
		sortOrder = "DESC"
	}

	query := fmt.Sprintf("SELECT * FROM projects p WHERE p.user_id = @user_id ORDER BY p.created %s", sortOrder)

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
		},
		PageSizeHint: int32(paginationOptions.PerPage),
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	projects := []Project{}
	currentPage := 0

	for pager.More() {

		if currentPage == paginationOptions.Page {
			response, err := pager.NextPage(ctx)
			if err != nil {
				return []Project{}, err
			}

			for _, bytes := range response.Items {
				project := Project{}
				err := json.Unmarshal(bytes, &project)
				if err != nil {
					return []Project{}, err
				}
				projects = append(projects, project)
			}

			return projects, nil

		} else {
			_, err := pager.NextPage(ctx)
			if err != nil {
				return []Project{}, err
			}
			currentPage++
		}

	}

	return projects, nil

}

func (env *env) UpsertProject(ctx context.Context, project Project) (Project, error) {

	env.logger.Info("Upserting project")

	container, err := env.client.NewContainer("projects")

	partitionKey := azcosmos.NewPartitionKeyString(project.UserID)

	marshalled, err := json.Marshal(project)
	if err != nil {
		return project, err
	}

	_, err = container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return project, err
	}

	return project, nil

}

func (env *env) RemoveProject(ctx context.Context, userID string, projectID string) (interface{}, error) {

	env.logger.Info("Removing project")

	container, err := env.client.NewContainer("projects")

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	response, err := container.DeleteItem(ctx, partitionKey, projectID, nil)
	if err != nil {
		return nil, err
	}

	return response, nil

}
