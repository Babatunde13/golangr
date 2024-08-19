package user

import "github.com/gin-gonic/gin"

func UserRoutes(userRouter *gin.RouterGroup) {
	controller := CreateUserController()
	userRouter.GET("/", controller.GetUsers())
	userRouter.POST("/", controller.CreateUser())
	userRouter.GET("/:id", controller.GetUserByID())
	userRouter.PUT("/:id", controller.UpdateUser())
	userRouter.DELETE("/:id", controller.DeleteUser())
	userRouter.POST("/login", controller.LoginUser())
	userRouter.GET("/me", controller.AuthMiddleware(), controller.AuthUser())
	userRouter.POST("/ai/suggest", controller.AuthMiddleware(), controller.SuggestWithGPT())
	userRouter.POST("/ai/tts", controller.AuthMiddleware(), controller.TextToSpeech())
	userRouter.POST("/upload", controller.AuthMiddleware(), controller.Upload())
	userRouter.GET("/download", controller.AuthMiddleware(), controller.Download())
}
