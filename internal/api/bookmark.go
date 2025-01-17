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
		c.JSON(401, err)
	}

	bookmarks, err := e.db.GetBookmarks(context.TODO(), cookie)
	if err != nil {
		c.JSON(400, err)
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
		c.JSON(401, err)
	}

	var req AddBookmarkReq
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(400, err)
	}

	bookmark := db.Bookmark{
		ID: 			   req.ID, //temporary
		Link: 			   req.Link,
		Label: 			   req.Label,
		UserID: 		   cookie,
	}

	//TODO: only add bookmark if it doesn't already exist for this user
	response, err := e.db.AddBookmark(context.TODO(), bookmark)
	if err != nil {
		c.JSON(400, err)
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
		c.JSON(401, err)
	}

	id := c.Param("bookmarkId")

	//TODO: verify user has this bookmark

	response, err := e.db.RemoveBookmark(context.TODO(), cookie, id)
	if err != nil {
		c.JSON(400, err)
	}

	c.JSON(204, response)

}