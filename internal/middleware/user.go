package middleware

import (
	"4h-recordbook-backend/pkg/db"
	"4h-recordbook-backend/internal/utils"
	"context"

	"github.com/gin-gonic/gin"
)

func GetUser(database db.Db) gin.HandlerFunc {
	return func(c* gin.Context){
		// Attempt to get the user

		_, err := database.GetUser(context.TODO(), c.GetString("user_id"));
		if(err != nil){
			// Force create user
			timestamp := utils.TimeNow()
			user := db.User{
				ID:                c.GetString("user_id"),
				Email:             c.GetString("user_email"),
				Birthdate:         "TODO",
				FirstName:         c.GetString("user_name"),
				MiddleNameInitial: "TODO",
				CountyName:        "TODO",
				GenericDatabaseInfo: db.GenericDatabaseInfo{
					Created: timestamp.String(),
					Updated: timestamp.String(),
				},
			}

			_, err := database.UpsertUser(context.TODO(), user)
			if err != nil {
				c.JSON(400, gin.H{
					"message": "Database Error >:(",
				})
			}
		}
		// User exists! Continue on.

		return;
	}
}
