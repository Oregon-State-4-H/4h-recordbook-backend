package db 

import (
	"time"
	"strconv"
	"context"
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type Project struct {
	ID string `json:"id"`
	Year string `json:"year"`
	Name string `json:"name"`
	Description string `json:"description"`
	Type string `json:"type"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
	UserID string `json:"userid"`
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

func (env *env) GetCurrentProjects(ctx context.Context, userID string) ([]Project, error) {

	env.logger.Info("Getting current projects")

	container, err := env.client.NewContainer("projects")
	if err != nil {
		return []Project{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	now := time.Now()
	year := strconv.Itoa(now.Year())
	
	query := "SELECT * FROM projects p WHERE p.userid = @user_id AND p.year = @year"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
			{Name: "@year", Value: year},
		},
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	projects := []Project{}

	for pager.More() {
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
	}

	return projects, nil

}

func (env *env) GetProjectsByUser(ctx context.Context, userID string) ([]Project, error) {

	env.logger.Info("Getting projects")

	container, err := env.client.NewContainer("projects")
	if err != nil {
		return []Project{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	query := "SELECT * FROM projects p WHERE p.userid = @user_id"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
		},
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	projects := []Project{}

	for pager.More() {
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
	}

	return projects, nil

}

func (env *env) UpsertProject(ctx context.Context, project Project) (interface{}, error) {
	
	env.logger.Info("Upserting project")

	container, err := env.client.NewContainer("projects")

	partitionKey := azcosmos.NewPartitionKeyString(project.UserID)

	marshalled, err := json.Marshal(project)
	if err != nil {
		return nil, err
	}

	response, err := container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return nil, err
	}

	return response, nil

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