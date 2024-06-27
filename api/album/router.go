package album

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func AlbumRouter(albumRouter *gin.RouterGroup){
	controller := CreateAlbumController()
	fmt.Println(("Creating a new API for our Album store"))
	albumRouter.GET("/", controller.GetAlbums)
	albumRouter.POST("/", controller.CreateAlbum)
	albumRouter.GET("/:id", controller.GetAlbumByID)
	albumRouter.PUT("/:id", controller.UpdateAlbum)
	albumRouter.DELETE("/:id", controller.DeleteAlbum)
}
