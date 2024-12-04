package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
type Bookmark struct {
	Link 				string		`json:"link" bson:"link"`
	Label 				string		`json:"label" bson:"label,omitempty"`
}
*/

type User struct {
	ID					primitive.ObjectID	`json:"id" bson:"_id"`
	Email 				string				`json:"email" bson:"email"`
	Birthdate 			time.Time			`json:"birthdate" bson:"birthdate"`
	First_name 			string				`json:"first_name" bson:"first_name"`
	Middle_name_initial string				`json:"middle_name_initial" bson:"middle_name_initial,omitempty"`
	Last_name_initial 	string				`json:"last_name_initial" bson:"last_name_initial"`
	County_name 		string				`json:"county_name" bson:"county_name"`
	Join_date 			time.Time			`json:"join_date" bson:"join_date"`
	//Bookmarks 			[]Bookmark	`json:"bookmarks" bson:"inline"`
}