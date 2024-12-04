package repository

import (
	"log"
	"context"
	. "4h-recordbook/backend/config"
	"4h-recordbook/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {}

var userConfig Config
var userCollection = new(mongo.Collection)
const userCollectionName = "users"

func init() {
	userConfig.Read()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(userConfig.Database.Uri))
	if err != nil {
		panic(err)
	}
	userCollection = client.Database(userConfig.Database.DatabaseName).Collection(userCollectionName)
}

func (u *UserRepository) Find(uid primitive.ObjectID) (models.User, error) {
	
	var user models.User
	filter := bson.D{{"_id", uid}}
	err := userCollection.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		log.Fatal(err)
	}

	return user, err

}

func (u *UserRepository) UpdateUser(user models.User) (interface{}, error) {
	
	filter := bson.D{{"_id", user.ID}}
	var update []bson.E
	update = append(update, bson.E{"birthdate", user.Birthdate})
	update = append(update, bson.E{"first_name", user.First_name})
	update = append(update, bson.E{"middle_name_initial", user.Middle_name_initial})
	update = append(update, bson.E{"last_name_initial", user.Last_name_initial})
	update = append(update, bson.E{"county_name", user.County_name})
	update = append(update, bson.E{"join_date", user.Join_date})
	
	result, err := userCollection.UpdateOne(context.TODO(), filter, bson.D{{"$set", update}})
	if err != nil {
		log.Fatal(err)
	}

	return result.UpsertedID, err

}

func (u *UserRepository) InsertUser(user models.User) (interface{}, error) {

	result, err := userCollection.InsertOne(context.TODO(), &user)
	if err != nil {
		log.Fatal(err)
	}

	return result.InsertedID, err

}