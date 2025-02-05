package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"4h-recordbook-backend/internal/utils"
	"4h-recordbook-backend/pkg/db"
	"github.com/beevik/guid"
)

type UpdateUserInput struct {
	Email 				string	`json:"email"`
	Birthdate 			string	`json:"birthdate"`
	FirstName 			string	`json:"first_name"`
	MiddleNameInitial 	string	`json:"middle_name_initial"`
	LastNameInitial 	string	`json:"last_name_initial"`
	CountyName 			string	`json:"county_name"`
}

type GetUserProfileOutput struct {
	User db.User `json:"user"`
}

// GetUserProfile godoc
// @Summary Get a user
// @Description Get user by JWT
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} api.GetUserProfileOutput 
// @Failure 401
// @Failure 404
// @Router /user [get]
func (e *env) getUserProfile(c *gin.Context) {

	//get id value from cookie, if no id value return 401
	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var output GetUserProfileOutput

	output.User, err = e.db.GetUser(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// UpdateUserProfile godoc
// @Summary Update a user
// @Description Update the signed-in user's information
// @Tags User
// @Accept json
// @Produce json
// @Param UpdateUserInput body api.UpdateUserInput true "User information"
// @Success 204 
// @Failure 400
// @Failure 401
// @Failure 404
// @Router /user [put]
func (e *env) updateUserProfile(c *gin.Context) {

	//get id value from cookie, if no id value return 401
	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var input UpdateUserInput
	err = c.BindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	user, err := e.db.GetUser(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedUser := db.User{
		ID: user.ID,
		Email: ternary(input.Email, user.Email),
		Birthdate: ternary(input.Birthdate, user.Birthdate),
		FirstName: ternary(input.FirstName, user.FirstName),
		MiddleNameInitial: ternary(input.MiddleNameInitial, user.MiddleNameInitial),
		CountyName: ternary(input.CountyName, user.CountyName),
		GenericDatabaseInfo: db.GenericDatabaseInfo {
			Created: user.Created,
			Updated: timestamp.ToString(),
		},
	}

	response, err := e.db.UpsertUser(context.TODO(), updatedUser)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}

//temporary signin and signup functions
type SignInInput struct {
	ID string `json:id" validate:"required"`
}

// SignIn godoc
// @Summary Sign in 
// @Description Placeholder route, sign in with ID
// @Tags User
// @Accept json
// @Produce json
// @Param ID body api.SignInInput true "User ID"
// @Success 204 
// @Failure 400
// @Failure 404
// @Router /signin [post]
func (e *env) signin(c *gin.Context) {

	var input SignInInput
	err := c.BindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	user, err := e.db.GetUser(context.TODO(), input.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.SetCookie("login_cookie", user.ID, 3600, "/", "localhost", false, false)
	c.JSON(204, nil)

}

// SignOut godoc
// @Summary Sign out
// @Description Placeholder route, sign out with JWT
// @Tags User
// @Accept json
// @Produce json
// @Success 204 
// @Failure 401
// @Router /signout [post]
func (e *env) signout(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	c.SetCookie("login_cookie", cookie, -1, "/", "localhost", false, false)
	c.JSON(204, nil)

}

type SignUpInput struct {
	Email 				string	`json:"email" validate:"required"`
	Birthdate 			string	`json:"birthdate" validate:"required"`
	FirstName 			string	`json:"first_name" validate:"required"`
	MiddleNameInitial 	string	`json:"middle_name_initial" validate:"required"`
	LastNameInitial 	string	`json:"last_name_initial" validate:"required"`
	CountyName 			string	`json:"county_name" validate:"required"`
}

// Signup godoc
// @Summary Sign up 
// @Description Placeholder route, sign up with custom user information
// @Tags User
// @Accept json
// @Produce json
// @Param ID body api.SignUpInput true "User information"
// @Success 204 
// @Failure 400
// @Failure 409
// @Router /signup [post]
func (e *env) signup(c *gin.Context) {

	var input SignUpInput
	err := c.BindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	g := guid.New()
	timestamp := utils.TimeNow()

	user := db.User{
		ID: 			   g.String(),
		Email: 			   input.Email,
		Birthdate: 		   input.Birthdate,
		FirstName: 		   input.FirstName,
		MiddleNameInitial: input.MiddleNameInitial,
		CountyName: 	   input.CountyName,
		GenericDatabaseInfo: db.GenericDatabaseInfo {
			Created: timestamp.ToString(),
			Updated: timestamp.ToString(),
		},
	}

	response, err := e.db.UpsertUser(context.TODO(), user)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}