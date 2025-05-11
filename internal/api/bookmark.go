package api

import (
	"4h-recordbook-backend/internal/utils"
	"4h-recordbook-backend/pkg/db"
	"context"
	"strconv"

	"github.com/beevik/guid"
	"github.com/gin-gonic/gin"
)

const BOOKMARKS = "bookmarks"

type GetBookmarksOutput struct {
	Bookmarks []db.Bookmark `json:"bookmarks"`
	Next      string        `json:"next"`
}

type GetBookmarkOutput struct {
	Bookmark db.Bookmark `json:"bookmark"`
}

type AddBookmarkInput struct {
	Link  string `json:"link" validate:"required"`
	Label string `json:"label" validate:"required"`
}

type AddBookmarkOutput GetBookmarkOutput

// GetUserBookmarks godoc
// @Summary Get all of a user's bookmarks
// @Description Returns an array of all the user's bookmarks, queried using JWT claims
// @Tags User Bookmarks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "Page number, default 0"
// @Param per_page query int false "Max number of items to return. Can be [1-100], default 30"
// @Param sort_by_newest query bool false "Sort order, default true"
// @Success 200 {object} api.GetBookmarksOutput
// @Failure 401
// @Router /bookmarks [get]
func (e *env) getUserBookmarks(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	pageStr := c.DefaultQuery("page", PAGE_DEFAULT_STR)
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(400, gin.H{
			"message": ErrQueryMustBeInt,
		})
		return
	}
	if page < 0 {
		page = PAGE_DEFAULT_INT
	}

	perPageStr := c.DefaultQuery("per_page", PER_PAGE_DEFAULT_STR)
	perPage, err := strconv.Atoi(perPageStr)
	if err != nil {
		c.JSON(400, gin.H{
			"message": ErrQueryMustBeInt,
		})
		return
	}
	if perPage < PER_PAGE_MIN_INT {
		perPage = PER_PAGE_MIN_INT
	} else if perPage > PER_PAGE_MAX_INT {
		perPage = PER_PAGE_MAX_INT
	}

	sortByNewestStr := c.DefaultQuery("sort_by_newest", "true")
	sortByNewest, err := strconv.ParseBool(sortByNewestStr)
	if err != nil {
		c.JSON(400, gin.H{
			"message": ErrQueryMustBeBool,
		})
		return
	}

	var output GetBookmarksOutput

	output.Bookmarks, err = e.db.GetBookmarks(context.TODO(), claims.ID, page, perPage, sortByNewest)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if len(output.Bookmarks) == perPage {
		output.Next = utils.NextUrl(c, BOOKMARKS, page, perPage, sortByNewest)
	}

	c.JSON(200, output)

}

// GetBookmarkByLink godoc
// @Summary Get a bookmark by the link
// @Description Returns a bookmark with the searched link, queried using JWT claims
// @Tags User Bookmarks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param link path string true "Bookmark link"
// @Success 200 {object} api.GetBookmarkOutput
// @Failure 401
// @Router /bookmarks/{link} [get]
func (e *env) getBookmarkByLink(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var output GetBookmarkOutput

	link := c.Param("link")

	output.Bookmark, err = e.db.GetBookmarkByLink(context.TODO(), claims.ID, link)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if output.Bookmark == (db.Bookmark{}) {
		c.JSON(404, gin.H{
			"message": ErrNotFound,
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
// @Success 201 {object} api.AddBookmarkOutput
// @Failure 400
// @Failure 401
// @Failure 409
// @Router /bookmarks [post]
func (e *env) addUserBookmark(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var input AddBookmarkInput
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

	bookmark := db.Bookmark{
		ID:     g.String(),
		Link:   input.Link,
		Label:  input.Label,
		UserID: claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: timestamp.String(),
			Updated: timestamp.String(),
		},
	}

	existingBookmark, err := e.db.GetBookmarkByLink(context.TODO(), claims.ID, input.Link)
	if existingBookmark != (db.Bookmark{}) {
		c.JSON(409, gin.H{
			"message": ErrBookmarkConflict,
		})
		return
	}
	if err != nil {
		response := InterpretCosmosError(err)
		if response.Code != 404 {
			c.JSON(response.Code, gin.H{
				"message": response.Message,
			})
			return
		}
	}

	var output AddBookmarkOutput

	output.Bookmark, err = e.db.AddBookmark(context.TODO(), bookmark)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(201, output)

}

// DeleteUserBookmark godoc
// @Summary Removes a bookmark
// @Description Deletes a user's bookmark given the bookmark ID
// @Tags User Bookmarks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param bookmarkID path string true "Bookmark ID"
// @Success 204
// @Failure 401
// @Failure 404
// @Router /bookmarks/{bookmarkID} [delete]
func (e *env) deleteUserBookmark(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	bookmarkID := c.Param("bookmarkID")

	response, err := e.db.RemoveBookmark(context.TODO(), claims.ID, bookmarkID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}
