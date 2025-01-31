package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"4h-recordbook-backend/internal/utils"
	"4h-recordbook-backend/pkg/db"
	"github.com/beevik/guid"
)

type UpsertProjectInput struct {
	Year 		string `json:"year" validate:"required"`
	Name 		string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Type 		string `json:"type" validate:"required"`
	StartDate   string `json:"start_date" validate:"required"`
	EndDate		string `json:"end_date" validate:"required"`
}

// GetCurrentProjects godoc
// @Summary 
// @Description 
// @Tags Projects
// @Accept json
// @Produce json
// @Success 200 
// @Router /projects [get]
func (e *env) getCurrentProjects(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	projects, err := e.db.GetCurrentProjects(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, projects)

}

// GetProjects godoc
// @Summary 
// @Description 
// @Tags Projects
// @Accept json
// @Produce json
// @Success 200 
// @Router /project [get]
func (e *env) getProjects(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	projects, err := e.db.GetProjectsByUser(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, projects)

}

// GetProject godoc
// @Summary 
// @Description 
// @Tags Projects
// @Accept json
// @Produce json
// @Success 200 
// @Router /project/{projectId} [get]
func (e *env) getProject(c *gin.Context) {
	
	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	id := c.Param("projectId")

	project, err := e.db.GetProjectByID(context.TODO(), cookie, id)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, project)

}

// AddProject godoc
// @Summary 
// @Description 
// @Tags Projects
// @Accept json
// @Produce json
// @Success 200 
// @Router /project [post]
func (e *env) addProject(c *gin.Context) {
	
	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var input UpsertProjectInput
	err = c.BindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	err = e.validator.Struct(input)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	g := guid.New()
	timestamp := utils.TimeNow()
	
	//verify StartDate and EndDate are properly formatted
	startDate, err := utils.StringToTimestamp(input.StartDate)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
	}

	endDate, err := utils.StringToTimestamp(input.EndDate)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
	}

	project := db.Project{
		ID: 			   g.String(),
		Year: 			   input.Year,
		Name: 			   input.Name,
		Description:	   input.Description,
		Type:			   input.Type,
		StartDate:		   startDate.ToString(),
		EndDate:		   endDate.ToString(),
		UserID: 		   cookie,
		Created:		   timestamp.ToString(),
		Updated:		   timestamp.ToString(),
	}

	response, err := e.db.UpsertProject(context.TODO(), project)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}

// UpdateProject godoc
// @Summary 
// @Description 
// @Tags Projects
// @Accept json
// @Produce json
// @Success 200 
// @Router /project/{projectId} [put]
func (e *env) updateProject(c *gin.Context) {
	
	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var input UpsertProjectInput
	err = c.BindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}
	
	err = e.validator.Struct(input)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	//verify StartDate and EndDate are properly formatted
	startDate, err := utils.StringToTimestamp(input.StartDate)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
	}

	endDate, err := utils.StringToTimestamp(input.EndDate)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
	}

	id := c.Param("projectId")

	project, err := e.db.GetProjectByID(context.TODO(), cookie, id)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedProject := db.Project{
		ID: 			   project.ID,
		Year: 			   input.Year,
		Name: 		   	   input.Name,
		Description: 	   input.Description,
		Type: 			   input.Type,
		StartDate: 	   	   startDate.ToString(),
		EndDate: 	       endDate.ToString(),
		UserID:			   cookie,
		Created: 		   project.Created,
		Updated:		   timestamp.ToString(),
	}

	response, err := e.db.UpsertProject(context.TODO(), updatedProject)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}