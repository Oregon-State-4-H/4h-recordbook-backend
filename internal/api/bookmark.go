package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"4h-recordbook-backend/internal/utils"
	"4h-recordbook-backend/pkg/db"
	"github.com/beevik/guid"
)

type AddBookmarkInput struct {
	Link string `json:"link" validate:"required"`
	Label string `json:"label" validate:"required"`
}

type GetBookmarksOutput struct {
	Bookmarks []db.Bookmark `json:"bookmarks"`
}

// GetUserBookmarks godoc
// @Summary Get all of a user's bookmarks
// @Description Returns an array of all the user's bookmarks, queried using JWT claims
// @Tags User Bookmarks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} api.GetBookmarksOutput
// @Failure 401 
// @Router /bookmarks [get]
func (e *env) getUserBookmarks(c *gin.Context) {
	
	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var output GetBookmarksOutput

	output.Bookmarks, err = e.db.GetBookmarks(context.TODO(), claims.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// AddUserBookmark godoc
// @Summary Adds a bookmark
// @Description Adds a bookmark to a user's personal records. 
// @Description The new bookmark can not have the same link as another of the user's bookmarks
// @Tags User Bookmarks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param AddBookmarkInput body api.AddBookmarkInput true "Bookmark information"
// @Success 204 
// @Failure 400
// @Failure 401
// @Failure 409
// @Router /bookmarks [post]
func (e *env) addUserBookmark(c *gin.Context) {
	
	claims, err := decodeJWT(c)
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
		ID: g.String(),
		Link: input.Link,
		Label: input.Label,
		UserID: claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo {
			Created: timestamp.ToString(),
			Updated: timestamp.ToString(),
		},
	}

	existingBookmark, err := e.db.GetBookmarkByLink(context.TODO(), claims.ID, input.Link)
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
// @Summary Removes a bookmark
// @Description Deletes a user's bookmark given the bookmark ID
// @Tags User Bookmarks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param bookmarkId path string true "Bookmark ID"
// @Success 204
// @Failure 401
// @Failure 404 
// @Router /bookmarks/{bookmarkId} [delete]
func (e *env) removeUserBookmark(c *gin.Context) {
	
	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	id := c.Param("bookmarkId")

	response, err := e.db.RemoveBookmark(context.TODO(), claims.ID, id)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}