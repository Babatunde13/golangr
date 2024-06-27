package album

import "github.com/gin-gonic/gin"

type IAlbumController interface {
	GetAlbums(c *gin.Context)
	CreateAlbum(c *gin.Context)
	GetAlbumByID(c *gin.Context)
	UpdateAlbum(c *gin.Context)
	DeleteAlbum(c *gin	.Context)
}

type AlbumController struct{}
