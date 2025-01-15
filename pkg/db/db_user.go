package db

import (
	"context"
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type User struct {
	ID					string	`json:"id"`
	Email 				string	`json:"email"`
	Birthdate 			string	`json:"birthdate"`
	FirstName 			string	`json:"first_name"`
	MiddleNameInitial 	string	`json:"middle_name_initial"`
	LastNameInitial 	string	`json:"last_name_initial"`
	CountyName 			string	`json:"county_name"`
	JoinDate 			string	`json:"join_date"`
}

func (env *env) GetUser(ctx context.Context, id string) (User, error) {
	
	env.logger.Info("Getting user")
	user := User{}

	container, err := env.client.NewContainer("users")
	if err != nil {
		return user, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(id)

	response, err := container.ReadItem(ctx, partitionKey, id, nil)
	if err != nil {
		return user, err
	}

	if response.RawResponse.StatusCode == 200 {
		err := json.Unmarshal(response.Value, &user)
		if err != nil {
			return user, err
		}
	}

	return user, nil

}

func (env *env) UpsertUser(ctx context.Context, user User) (interface{}, error) {

	env.logger.Info("Upserting user")

	container, err := env.client.NewContainer("users")

	partitionKey := azcosmos.NewPartitionKeyString(user.ID)

	marshalled, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	response, err := container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return nil, err
	}

	return response, nil

}