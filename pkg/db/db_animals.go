package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type Animal struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Species         string  `json:"species"`
	BirthDate       string  `json:"birth_date"`
	PurchaseDate    string  `json:"purchase_date"`
	SireBreed       string  `json:"sire_breed"`
	DamBreed        string  `json:"dam_breed"`
	BeginningWeight float64 `json:"beginning_weight"`
	BeginningDate   string  `json:"beginning_date"`
	EndWeight       float64 `json:"end_weight"`
	EndDate         string  `json:"end_date"`
	AnimalCost      float64 `json:"animal_cost"`
	SalePrice       float64 `json:"sale_price"`
	YieldGrade      string  `json:"yield_grade"`
	QualityGrade    string  `json:"quality_grade"`
	UserID          string  `json:"user_id"`
	ProjectID       string  `json:"project_id"`
	GenericDatabaseInfo
}

func (a Animal) GetID() string {
	return a.ID
}

func (env *env) GetAnimalsByProject(ctx context.Context, userID string, projectID string, paginationOptions PaginationOptions) ([]Animal, error) {

	env.logger.Info("Getting animals by project")

	container, err := env.client.NewContainer("animals")
	if err != nil {
		return []Animal{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	sortOrder := "ASC"
	if paginationOptions.SortByNewest {
		sortOrder = "DESC"
	}

	query := fmt.Sprintf("SELECT * FROM animals a WHERE a.user_id = @user_id AND a.project_id = @project_id ORDER BY a.created %s", sortOrder)

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
			{Name: "@project_id", Value: projectID},
		},
		PageSizeHint: int32(paginationOptions.PerPage),
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	animals := []Animal{}
	currentPage := 0

	for pager.More() {

		if currentPage == paginationOptions.Page {
			response, err := pager.NextPage(ctx)
			if err != nil {
				return []Animal{}, err
			}

			for _, bytes := range response.Items {
				animal := Animal{}
				err := json.Unmarshal(bytes, &animal)
				if err != nil {
					return []Animal{}, err
				}
				animals = append(animals, animal)
			}

			return animals, nil

		} else {
			_, err := pager.NextPage(ctx)
			if err != nil {
				return []Animal{}, err
			}
			currentPage++
		}

	}

	return animals, nil

}

func (env *env) GetProjectDependentAnimals(ctx context.Context, userID string, projectID string) ([]Identifiable, error) {

	env.logger.Info("Getting project dependent animals")

	container, err := env.client.NewContainer("animals")
	if err != nil {
		return []Identifiable{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	query := "SELECT * FROM animals a WHERE a.user_id = @user_id AND a.project_id = @project_id"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
			{Name: "@project_id", Value: projectID},
		},
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	animals := []Animal{}

	for pager.More() {

		response, err := pager.NextPage(ctx)
		if err != nil {
			return []Identifiable{}, err
		}

		for _, bytes := range response.Items {
			animal := Animal{}
			err := json.Unmarshal(bytes, &animal)
			if err != nil {
				return []Identifiable{}, err
			}
			animals = append(animals, animal)
		}

	}

	identifiables := []Identifiable{}

	for _, a := range animals {
		identifiables = append(identifiables, a)
	}

	return identifiables, nil

}

func (env *env) GetAnimalByID(ctx context.Context, userID string, animalID string) (Animal, error) {

	env.logger.Info("Getting animal by ID")
	animal := Animal{}

	container, err := env.client.NewContainer("animals")
	if err != nil {
		return animal, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	response, err := container.ReadItem(ctx, partitionKey, animalID, nil)
	if err != nil {
		return animal, err
	}

	err = json.Unmarshal(response.Value, &animal)
	if err != nil {
		return animal, err
	}

	return animal, nil

}

func (env *env) UpsertAnimal(ctx context.Context, animal Animal) (Animal, error) {

	env.logger.Info("Upserting animal")

	container, err := env.client.NewContainer("animals")
	if err != nil {
		return animal, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(animal.UserID)

	marshalled, err := json.Marshal(animal)
	if err != nil {
		return animal, err
	}

	_, err = container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return animal, err
	}

	return animal, nil

}

func (env *env) RemoveAnimal(ctx context.Context, userID string, animalID string) (interface{}, error) {

	env.logger.Info("Removing animal")

	container, err := env.client.NewContainer("animals")
	if err != nil {
		return nil, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	response, err := container.DeleteItem(ctx, partitionKey, animalID, nil)
	if err != nil {
		return nil, err
	}

	for _, dependent := range env.dependentsMap["animals"] {
		identifiables, err := dependent.GetRelated(ctx, userID, animalID)
		if err != nil {
			return nil, err
		}
		for _, identifiable := range identifiables {
			_, err := dependent.Delete(ctx, userID, identifiable.GetID())
			if err != nil {
				return nil, err
			}
		}
	}

	return response, nil

}
