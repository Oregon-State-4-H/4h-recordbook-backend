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

type ProjectRepository struct {}

var projectConfig Config
var projectCollection = new(mongo.Collection)
const projectCollectionName = "projects"

func init() {
	projectConfig.Read()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(projectConfig.Database.Uri))
	if err != nil {
		panic(err)
	}
	projectCollection = client.Database(projectConfig.Database.DatabaseName).Collection(projectCollectionName)
}

func (p *ProjectRepository) FindUsersProjects(uid primitive.ObjectID) ([]models.Project, error) {

	var projects []models.Project

	filter := bson.D{{"_uid", uid}}

	cursor, err := projectCollection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	if err = cursor.All(context.TODO(), &projects); err != nil {
		log.Fatal(err)
	}

	return projects, err

}

func (p *ProjectRepository) InsertProject(project models.Project) (interface{}, error){

	result, err := projectCollection.InsertOne(context.TODO(), &project)
	if err != nil {
		log.Fatal(err)
	}

	return result.InsertedID, err

}