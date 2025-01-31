package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"4h-recordbook-backend/internal/utils"
	"4h-recordbook-backend/pkg/db"
	"github.com/beevik/guid"
)

type AddBookmarkInput struct {
	Link  string `json:"link" validate:"required"`
	Label string `json:"label" validate:"required"`
}

// GetUserBookmarks godoc
// @Summary 
// @Description 
// @Tags User Bookmarks
// @Accept json
// @Produce json
// @Success 200 
// @Router /bookmarks [get]
func (e *env) getUserBookmarks(c *gin.Context) {
	
	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	bookmarks, err := e.db.GetBookmarks(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, bookmarks)

}

// AddUserBookmark godoc
// @Summary 
// @Description 
// @Tags User Bookmarks
// @Accept json
// @Produce json
// @Success 200 
// @Router /bookmarks [post]
func (e *env) addUserBookmark(c *gin.Context) {
	
	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var input AddBookmarkInput
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

	bookmark := db.Bookmark{
		ID: 			   g.String(),
		Link: 			   input.Link,
		Label: 			   input.Label,
		UserID: 		   cookie,
		Created:		   timestamp.ToString(),
		Updated:		   timestamp.ToString(),
	}

	existingBookmark, err := e.db.GetBookmarkByLink(context.TODO(), cookie, input.Link)
	if existingBookmark != (db.Bookmark{}) {
		c.JSON(409, gin.H{
			"message": HTTPResponseCodeMap[409],
		})
		return
	}

	response, err := e.db.AddBookmark(context.TODO(), bookmark)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}

// RemoveUserBookmark godoc
// @Summary 
// @Description 
// @Tags User Bookmarks
// @Accept json
// @Produce json
// @Success 200 
// @Router /bookmarks/{bookmarkId} [delete]
func (e *env) removeUserBookmark(c *gin.Context) {
	
	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	id := c.Param("bookmarkId")

	response, err := e.db.RemoveBookmark(context.TODO(), cookie, id)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}