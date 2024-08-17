package api

import (
	"github.com/Babatunde13/golangr/api/album"
	"github.com/Babatunde13/golangr/api/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

const DocsURL = "https://documenter.getpostman.com/view/13469769/TzJx9G7m"

func DocsController (c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, DocsURL)
}

func V1Routes (r *gin.RouterGroup) {
	r.GET("/docs", DocsController)
	album.AlbumRouter(r.Group("/albums"))
	user.UserRoutes(r.Group("/users"))
}

func HomeController (c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welcome to our API",
	})
}

func Router (r *gin.Engine) {
	apiV1Router := r.Group("/api/v1")
	apiV1Router.GET("/", HomeController)
	V1Routes(apiV1Router)
}
