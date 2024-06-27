package user

import "github.com/gin-gonic/gin"

type IUserController interface {
	GetUsers(c *gin.Context)
	CreateUser(c *gin.Context)
	GetUserByID(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type UserController struct{}

