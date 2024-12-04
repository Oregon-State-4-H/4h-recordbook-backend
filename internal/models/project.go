package models

import (
	//"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Project struct {
	//ID			primitive.ObjectID `json:"_id" bson:"_id"`
	Year 		string	`json:"year" bson:"year"`
	ProjectName string	`json:"projectname" bson:"projectname"`
	Description string	`json:"description" bson:"description,omitempty`
	UserID		primitive.ObjectID `json:"uid" bson:"_uid"`
	/*
	Type String
	StartDate time.Time
	EndDate time.Time
	*/
}