package api

import (
	"4h-recordbook-backend/internal/utils"
	"4h-recordbook-backend/pkg/db"
	"strconv"

	"github.com/beevik/guid"
	"github.com/gin-gonic/gin"
)

type GetProjectsOutput struct {
	Projects []db.Project `json:"projects"`
	Next     string       `json:"next"`
}

type GetProjectOutput struct {
	Project db.Project `json:"project"`
}

type UpsertProjectInput struct {
	Year        string `json:"year" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Type        string `json:"type" validate:"required"`
	StartDate   string `json:"start_date" validate:"required"`
	EndDate     string `json:"end_date" validate:"required"`
}

type UpsertProjectOutput GetProjectOutput

// GetCurrentProjects godoc
// @Summary Gets projects of the current year
// @Description Gets all of a user's projects that take place in the last 12 months
// @Tags Project
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "Page number, default 0"
// @Param per_page query int false "Max number of items to return. Can be [1-200], default 100"
// @Param sort_by_newest query bool false "Sort results by most recently added, default true"
// @Success 200 {object} api.GetProjectsOutput
// @Failure 401
// @Router /projects [get]
func (e *env) getCurrentProjects(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var output GetProjectsOutput

	paginationOptions := db.PaginationOptions{
		Page:         c.GetInt(CONTEXT_KEY_PAGE),
		PerPage:      c.GetInt(CONTEXT_KEY_PER_PAGE),
		SortByNewest: c.GetBool(CONTEXT_KEY_SORT_BY_NEWEST),
	}

	output.Projects, err = e.db.GetCurrentProjects(c.Request.Context(), claims.ID, paginationOptions)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if len(output.Projects) == paginationOptions.PerPage {

		queryParamsMap := make(map[string]string)
		queryParamsMap[CONTEXT_KEY_PAGE] = strconv.Itoa(paginationOptions.Page + 1)
		queryParamsMap[CONTEXT_KEY_PER_PAGE] = strconv.Itoa(paginationOptions.PerPage)
		queryParamsMap[CONTEXT_KEY_SORT_BY_NEWEST] = strconv.FormatBool(paginationOptions.SortByNewest)

		nextUrlInput := utils.NextUrlInput{
			Context:     c,
			QueryParams: queryParamsMap,
		}

		output.Next = utils.BuildNextUrl(nextUrlInput)
	}

	c.JSON(200, output)

}

// GetProjects godoc
// @Summary Get all of a user's projects
// @Description Gets all of a user's saved projects regardless of year
// @Tags Project
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "Page number, default 0"
// @Param per_page query int false "Max number of items to return. Can be [1-200], default 100"
// @Param sort_by_newest query bool false "Sort results by most recently added, default true"
// @Success 200 {object} api.GetProjectsOutput
// @Failure 401
// @Router /project [get]
func (e *env) getProjects(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var output GetProjectsOutput

	paginationOptions := db.PaginationOptions{
		Page:         c.GetInt(CONTEXT_KEY_PAGE),
		PerPage:      c.GetInt(CONTEXT_KEY_PER_PAGE),
		SortByNewest: c.GetBool(CONTEXT_KEY_SORT_BY_NEWEST),
	}

	output.Projects, err = e.db.GetProjectsByUser(c.Request.Context(), claims.ID, paginationOptions)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if len(output.Projects) == paginationOptions.PerPage {

		queryParamsMap := make(map[string]string)
		queryParamsMap[CONTEXT_KEY_PAGE] = strconv.Itoa(paginationOptions.Page + 1)
		queryParamsMap[CONTEXT_KEY_PER_PAGE] = strconv.Itoa(paginationOptions.PerPage)
		queryParamsMap[CONTEXT_KEY_SORT_BY_NEWEST] = strconv.FormatBool(paginationOptions.SortByNewest)

		nextUrlInput := utils.NextUrlInput{
			Context:     c,
			QueryParams: queryParamsMap,
		}

		output.Next = utils.BuildNextUrl(nextUrlInput)
	}

	c.JSON(200, output)

}

// GetProject godoc
// @Summary Get a project
// @Description Get a user's project by ID
// @Tags Project
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param projectID path string true "Project ID"
// @Success 200 {object} api.GetProjectOutput
// @Failure 401
// @Failure 404
// @Router /project/{projectID} [get]
func (e *env) getProject(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	projectID := c.Param("projectID")

	var output GetProjectOutput

	output.Project, err = e.db.GetProjectByID(c.Request.Context(), claims.ID, projectID)
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
// @Tags Project
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param UpsertProjectInput body api.UpsertProjectInput true "Project information"
// @Success 201 {object} api.UpsertProjectOutput
// @Failure 400
// @Failure 401
// @Router /project [post]
func (e *env) addProject(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var input UpsertProjectInput
	err = c.BindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{
			"message": ErrBadRequest,
		})
		return
	}

	err = e.validator.Struct(input)
	if err != nil {
		c.JSON(400, gin.H{
			"message": ErrMissingFields,
		})
		return
	}

	g := guid.New()
	timestamp := utils.TimeNow()

	startDate, err := utils.StringToTimestamp(input.StartDate)
	if err != nil {
		c.JSON(400, gin.H{
			"message": ErrBadDate,
		})
		return
	}

	endDate, err := utils.StringToTimestamp(input.EndDate)
	if err != nil {
		c.JSON(400, gin.H{
			"message": ErrBadDate,
		})
		return
	}

	project := db.Project{
		ID:          g.String(),
		Year:        input.Year,
		Name:        input.Name,
		Description: input.Description,
		Type:        input.Type,
		StartDate:   startDate.String(),
		EndDate:     endDate.String(),
		UserID:      claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: timestamp.String(),
			Updated: timestamp.String(),
		},
	}

	var output UpsertProjectOutput

	output.Project, err = e.db.UpsertProject(c.Request.Context(), project)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(201, output)

}

// UpdateProject godoc
// @Summary Update a project
// @Description Updates a user's project information
// @Tags Project
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param projectID path string true "Project ID"
// @Param UpsertProjectInput body api.UpsertProjectInput true "Project information"
// @Success 200 {object} api.UpsertProjectOutput
// @Failure 400
// @Failure 401
// @Failure 404
// @Router /project/{projectID} [put]
func (e *env) updateProject(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var input UpsertProjectInput
	err = c.BindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{
			"message": ErrBadRequest,
		})
		return
	}

	err = e.validator.Struct(input)
	if err != nil {
		c.JSON(400, gin.H{
			"message": ErrMissingFields,
		})
		return
	}

	startDate, err := utils.StringToTimestamp(input.StartDate)
	if err != nil {
		c.JSON(400, gin.H{
			"message": ErrBadDate,
		})
		return
	}

	endDate, err := utils.StringToTimestamp(input.EndDate)
	if err != nil {
		c.JSON(400, gin.H{
			"message": ErrBadDate,
		})
		return
	}

	projectID := c.Param("projectID")

	project, err := e.db.GetProjectByID(c.Request.Context(), claims.ID, projectID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedProject := db.Project{
		ID:          project.ID,
		Year:        input.Year,
		Name:        input.Name,
		Description: input.Description,
		Type:        input.Type,
		StartDate:   startDate.String(),
		EndDate:     endDate.String(),
		UserID:      claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: project.Created,
			Updated: timestamp.String(),
		},
	}

	var output UpsertProjectOutput

	output.Project, err = e.db.UpsertProject(c.Request.Context(), updatedProject)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// DeleteProject godoc
// @Summary Removes a project
// @Description Deletes a user's project given the project ID
// @Tags Project
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param projectID path string true "Project ID"
// @Success 204
// @Failure 401
// @Failure 404
// @Router /project/{projectID} [delete]
func (e *env) deleteProject(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	projectID := c.Param("projectID")

	response, err := e.db.RemoveProject(c.Request.Context(), claims.ID, projectID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}
