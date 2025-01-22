package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"4h-recordbook-backend/internal/utils"
	"4h-recordbook-backend/pkg/db"
)

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

	if project == (db.Project{}){
		c.JSON(404, gin.H{
			"message": HTTPResponseCodeMap[404],
		})
		return
	}

	c.JSON(200, project)

}

type AddProjectReq struct {
	ID			string `json:"id"`
	Year 		string `json:"year"`
	Name 		string `json:"name"`
	Description string `json:"description"`
	Type 		string `json:"type"`
	StartDate   string `json:"start_date"`
	EndDate		string `json:"end_date"`
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

	var req AddProjectReq
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(500, gin.H{
			"message": HTTPResponseCodeMap[500],
		})
		return
	}

	timestamp := utils.TimeNow()
	
	//verify StartDate and EndDate are properly formatted
	startDate, err := utils.StringToTimestamp(req.StartDate)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
	}

	endDate, err := utils.StringToTimestamp(req.EndDate)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
	}

	project := db.Project{
		ID: 			   req.ID, //temporary
		Year: 			   req.Year,
		Name: 			   req.Name,
		Description:	   req.Description,
		Type:			   req.Type,
		StartDate:		   startDate.ToString(),
		EndDate:		   endDate.ToString(),
		UserID: 		   cookie,
		Created:		   timestamp.ToString(),
		Updated:		   timestamp.ToString(),
	}

	existingProject, err := e.db.GetProjectByID(context.TODO(), cookie, req.ID)
	if existingProject != (db.Project{}) {
		c.JSON(409, gin.H{
			"message": HTTPResponseCodeMap[409],
		})
		return
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

type UpdateProjectReq struct {
	Year 		string `json:"year"`
	Name 		string `json:"name"`
	Description string `json:"description"`
	Type 		string `json:"type"`
	StartDate   string `json:"start_date"`
	EndDate		string `json:"end_date"`
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

	var req UpdateProjectReq
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(500, gin.H{
			"message": HTTPResponseCodeMap[500],
		})
		return
	}

	//verify StartDate and EndDate are properly formatted
	if req.StartDate != "" {
		_, err = utils.StringToTimestamp(req.StartDate)
		if err != nil {
			c.JSON(400, gin.H{
				"message": HTTPResponseCodeMap[400],
			})
		return
		}
	}

	if req.EndDate != "" {
		_, err = utils.StringToTimestamp(req.EndDate)
		if err != nil {
			c.JSON(400, gin.H{
				"message": HTTPResponseCodeMap[400],
			})
			return
		}
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

	if project == (db.Project{}){
		c.JSON(404, gin.H{
			"message": HTTPResponseCodeMap[404],
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedProject := db.Project{
		ID: 			   project.ID,
		Year: 			   ternary(req.Year, project.Year),
		Name: 		   	   ternary(req.Name, project.Name),
		Description: 	   ternary(req.Description, project.Description),
		Type: 			   ternary(req.Type, project.Type),
		StartDate: 	   	   ternary(req.StartDate, project.StartDate),
		EndDate: 	       ternary(req.EndDate, project.EndDate),
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