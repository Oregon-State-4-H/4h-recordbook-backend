package middleware

import (
	"4h-recordbook-backend/pkg/db"
	"context"

	"github.com/gin-gonic/gin"
)

func GetUser(database db.Db) gin.HandlerFunc {
	return func(c* gin.Context){
		// Attempt to get the user

		_, err := database.GetUser(context.TODO(), c.GetString("user_id"));
		if(err != nil){
			c.JSON(422, gin.H{
				"error": "not_registered",
			})
			return;
		}
		// User exists! Continue on.

		return;
	}
}
