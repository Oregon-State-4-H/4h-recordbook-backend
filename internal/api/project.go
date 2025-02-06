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

type GetProjectsOutput struct {
	Projects []db.Project `json:"projects"`
}

type GetProjectOutput struct {
	Project db.Project `json:"project"`
}

// GetCurrentProjects godoc
// @Summary Gets projects of the current year
// @Description Gets all of a user's projects that take place in the last 12 months
// @Tags Projects
// @Accept json 
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} api.GetProjectsOutput
// @Failure 401
// @Router /projects [get]
func (e *env) getCurrentProjects(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var output GetProjectsOutput

	output.Projects, err = e.db.GetCurrentProjects(context.TODO(), claims.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// GetProjects godoc
// @Summary Get all of a user's projects
// @Description Gets all of a user's saved projects regardless of year
// @Tags Projects
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} api.GetProjectsOutput
// @Failure 401
// @Router /project [get]
func (e *env) getProjects(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var output GetProjectsOutput

	output.Projects, err = e.db.GetProjectsByUser(context.TODO(), claims.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// GetProject godoc
// @Summary Get a project
// @Description Get a user's project by ID
// @Tags Projects
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param projectId path string true "Project ID"
// @Success 200 {object} api.GetProjectOutput
// @Failure 401
// @Failure 404
// @Router /project/{projectId} [get]
func (e *env) getProject(c *gin.Context) {
	
	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	id := c.Param("projectId")

	var output GetProjectOutput

	output.Project, err = e.db.GetProjectByID(context.TODO(), claims.ID, id)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// AddProject godoc
// @Summary Add a project
// @Description Adds a project to a user's personal records
// @Tags Projects
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param UpsertProjectInput body api.UpsertProjectInput true "Project information"
// @Success 204
// @Failure 400
// @Failure 401
// @Router /project [post]
func (e *env) addProject(c *gin.Context) {
	
	claims, err := decodeJWT(c)
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
		return
	}

	endDate, err := utils.StringToTimestamp(input.EndDate)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	project := db.Project{
		ID: 			   g.String(),
		Year: 			   input.Year,
		Name: 			   input.Name,
		Description:	   input.Description,
		Type:			   input.Type,
		StartDate:		   startDate.ToString(),
		EndDate:		   endDate.ToString(),
		UserID: 		   claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo {
			Created: timestamp.ToString(),
			Updated: timestamp.ToString(),
		},
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
// @Summary Update a project
// @Description Updates a user's project information
// @Tags Projects
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param projectId path string true "Project ID"
// @Param UpsertProjectInput body api.UpsertProjectInput true "Project information"
// @Success 204 
// @Failure 400
// @Failure 401
// @Failure 404
// @Router /project/{projectId} [put]
func (e *env) updateProject(c *gin.Context) {
	
	claims, err := decodeJWT(c)
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
		return
	}

	endDate, err := utils.StringToTimestamp(input.EndDate)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	id := c.Param("projectId")

	project, err := e.db.GetProjectByID(context.TODO(), claims.ID, id)
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
		UserID:			   claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo {
			Created: project.Created,
			Updated: timestamp.ToString(),
		},
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