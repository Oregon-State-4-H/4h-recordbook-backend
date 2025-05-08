package middleware

import (
	"4h-recordbook-backend/pkg/db"
	"context"

	"github.com/gin-gonic/gin"
)


func createUser(e *env, c *gin.Context) {
	e.logger.Info("Creating account for " + c.GetString("user_name"))
	e.logger.Warn("Required updating of Middle Name, Last Name, Birthdate, and County!")

	userInfo = User{
		ID: c.GetString("user_id"),
		Email: c.GetString("user_email"),
		Birthdate: "0-0-1970",
		FirstName: c.GetString("user_name"),
		MiddleNameInitial: "",
		LastNameInitial: "",
		CountyName: "",
	}

	response, err := e.db.UpsertUser(context.TODO(), userInfo)
	if(err){
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}
}

func (e *env) GetUser() gin.HandlerFunc {
	return func(c* gin.Context){
		// Attempt to get the user
		var user db.User;

		user, err := e.db.GetUser(context.TODO(), c.GetString("user_id"));
		if(err){
			// User might not yet exist
			createUser(e, c);
			return;
		}
		// User exists! Continue on.
		return;
	}
}
