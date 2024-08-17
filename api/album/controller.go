package album

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Babatunde13/golangr/api/database"
	"github.com/Babatunde13/golangr/api/response"
)

func CreateAlbumController() *AlbumController {
	var albumCtrl IAlbumController = &AlbumController{}
	return albumCtrl.(*AlbumController)
}

var albumColl = database.AlbumCollection()

func (ac *AlbumController) GetAlbums(c *gin.Context) {
	albums, err := albumColl.GetAllAlbums(database.Album{}); if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err, "Failed to retrieve all albums"))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse(albums, "Successfully retrieved all albums"))
}

func (ac *AlbumController) CreateAlbum(c *gin.Context) {
	var album database.Album
	if err := c.BindJSON(&album); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err, "Failed to create album"))
		return
	}

	newAlbum, err := albumColl.CreateAlbum(album); if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err, "Failed to create album"))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse(newAlbum, "Successfully created a new album"))
}

func (ac *AlbumController) GetAlbumByID(c *gin.Context) {
	id := c.Param("id")
	album, err := albumColl.GetAlbumByID(id); if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(err, "Album not found"))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(album, "Successfully retrieved album by ID"))
}

func (ac *AlbumController) UpdateAlbum(c *gin.Context) {
	id := c.Param("id")
	var album database.Album
	if err := c.BindJSON(&album); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err, "Failed to update album"))
		return
	}

	if updatedAlbum, err := albumColl.UpdateAlbum(id, album); err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(err, "Album not found"))
		return
	} else {
		c.JSON(http.StatusOK, response.SuccessResponse(updatedAlbum, "Successfully updated album"))
	}
}

func (ac *AlbumController) DeleteAlbum(c *gin.Context) {
	id := c.Param("id")
	album, err := albumColl.DeleteAlbum(id); if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(err, "Album not found"))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(album, "Successfully deleted album"))
}
