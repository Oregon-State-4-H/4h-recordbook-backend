package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"4h-recordbook-backend/pkg/db"
)

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

type AddBookmarkReq struct {
	ID	  string
	Link  string
	Label string
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

	var req AddBookmarkReq
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(500, gin.H{
			"message": HTTPResponseCodeMap[500],
		})
		return
	}

	bookmark := db.Bookmark{
		ID: 			   req.ID, //temporary
		Link: 			   req.Link,
		Label: 			   req.Label,
		UserID: 		   cookie,
	}

	e.logger.Info("1")

	existingBookmark, err := e.db.GetBookmarkByLink(context.TODO(), cookie, req.Link)
	if existingBookmark != (db.Bookmark{}) {
		e.logger.Info(existingBookmark)
		c.JSON(409, gin.H{
			"message": HTTPResponseCodeMap[409],
		})
		return
	}

	e.logger.Info("2")

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