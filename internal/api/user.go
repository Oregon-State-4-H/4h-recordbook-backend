package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"4h-recordbook-backend/pkg/db"
)

func ternary(s1 string, s2 string) (string){
	if s1 == "" {
		return s2
	}
	return s1
}

// GetUserProfile godoc
// @Summary 
// @Description 
// @Tags User
// @Accept json
// @Produce json
// @Success 200 
// @Router /user [get]
func (e *env) getUserProfile(c *gin.Context) {

	//get id value from cookie, if no id value return 401
	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, err)
	}

	user, err := e.db.GetUser(context.TODO(), cookie)
	if err != nil {
		c.JSON(400, err)
		return
	}

	c.JSON(200, user)

}

type UpdateUserReq struct {
	Email 				string	`json:"email"`
	Birthdate 			string	`json:"birthdate"`
	FirstName 			string	`json:"first_name"`
	MiddleNameInitial 	string	`json:"middle_name_initial"`
	LastNameInitial 	string	`json:"last_name_initial"`
	CountyName 			string	`json:"county_name"`
}

// UpdateUserProfile godoc
// @Summary 
// @Description 
// @Tags User
// @Accept json
// @Produce json
// @Success 200 
// @Router /user [put]
func (e *env) updateUserProfile(c *gin.Context) {

	//get id value from cookie, if no id value return 401
	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, err)
	}

	var req UpdateUserReq
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(400, err)
	}

	user, err := e.db.GetUser(context.TODO(), cookie)
	if err != nil {
		c.JSON(400, err)
		return
	}

	updatedUser := db.User{
		ID: 			   user.ID,
		Email: 			   ternary(req.Email, user.Email),
		Birthdate: 		   ternary(req.Birthdate, user.Birthdate),
		FirstName: 		   ternary(req.FirstName, user.FirstName),
		MiddleNameInitial: ternary(req.MiddleNameInitial, user.MiddleNameInitial),
		CountyName: 	   ternary(req.CountyName, user.CountyName),
		JoinDate: 		   user.JoinDate,
	}

	response, err := e.db.UpsertUser(context.TODO(), updatedUser)
	if err != nil {
		c.JSON(400, err)
	}

	c.JSON(204, response)

}

//temporary signin and signup functions
type SignInReq struct {
	ID string `json:id" validate:"required"`
}


func (e *env) signin(c *gin.Context) {

	var req SignInReq
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, err)
	}

	user, err := e.db.GetUser(context.TODO(), req.ID)
	if err != nil {
		c.JSON(400, err)
		return
	}

	c.SetCookie("login_cookie", user.ID, 3600, "/", "localhost", false, false)
	c.JSON(204, nil)

}

func (e *env) signout(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(400, err)
	}

	c.SetCookie("login_cookie", cookie, -1, "/", "localhost", false, false)
	c.JSON(204, nil)

}

type SignUpReq struct {
	ID					string  `json:"id" validate:"required"`
	Email 				string	`json:"email" validate:"required"`
	Birthdate 			string	`json:"birthdate" validate:"required"`
	FirstName 			string	`json:"first_name" validate:"required"`
	MiddleNameInitial 	string	`json:"middle_name_initial" validate:"required"`
	LastNameInitial 	string	`json:"last_name_initial" validate:"required"`
	CountyName 			string	`json:"county_name" validate:"required"`
}

func (e *env) signup(c *gin.Context) {

	var req SignUpReq
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, err)
	}

	user := db.User{
		ID: 			   req.ID,
		Email: 			   req.Email,
		Birthdate: 		   req.Birthdate,
		FirstName: 		   req.FirstName,
		MiddleNameInitial: req.MiddleNameInitial,
		CountyName: 	   req.CountyName,
		JoinDate: 		   "TODAY",
	}

	response, err := e.db.UpsertUser(context.TODO(), user)
	if err != nil {
		c.JSON(400, err)
	}

	c.JSON(204, response)

}