package handlers

import (
	"github.com/gin-gonic/gin"
	. "4h-recordbook/backend/internal/repository"
)

var projectRepository ProjectRepository

// GetCurrentProjects godoc
// @Summary 
// @Description 
// @Tags Projects
// @Accept json
// @Produce json
// @Success 200 
// @Router /projects [get]
func GetCurrentProjects(c *gin.Context) {

}

// GetProjects godoc
// @Summary 
// @Description 
// @Tags Projects
// @Accept json
// @Produce json
// @Success 200 
// @Router /project [get]
func GetProjects(c *gin.Context) {
	
}

// GetProject godoc
// @Summary 
// @Description 
// @Tags Projects
// @Accept json
// @Produce json
// @Success 200 
// @Router /project/{projectId} [get]
func GetProject(c *gin.Context) {
	
}

// AddProject godoc
// @Summary 
// @Description 
// @Tags Projects
// @Accept json
// @Produce json
// @Success 200 
// @Router /project [post]
func AddProject(c *gin.Context) {
	
}

// UpdateProject godoc
// @Summary 
// @Description 
// @Tags Projects
// @Accept json
// @Produce json
// @Success 200 
// @Router /project/{projectId} [put]
func UpdateProject(c *gin.Context) {
	
}