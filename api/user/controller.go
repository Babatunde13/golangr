package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Babatunde13/golangr/api/config"
	"github.com/Babatunde13/golangr/api/database"
	internalhttp "github.com/Babatunde13/golangr/api/http"
	"github.com/Babatunde13/golangr/api/response"
	"github.com/Babatunde13/golangr/api/utils"
)

func CreateUserController() IUserController {
	userCtrl := &UserController{}
	return IUserController(userCtrl)
}

var userColl = database.UserCollection()

func (uc *UserController) GetUsers() gin.HandlerFunc {
	return func (c *gin.Context) {
		users, err := userColl.GetAllUsers(database.User{}); if err != nil {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse(err, "Failed to retrieve all users"))
			return
		}

		c.JSON(http.StatusOK, response.SuccessResponse(users, "Successfully retrieved all users"))
	}
}

func (uc *UserController) CreateUser() gin.HandlerFunc {
	return func (c *gin.Context) {
		var user database.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(err, err.Error()))
			return
		}

		newUser, err := userColl.CreateUser(user); if err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(err, err.Error()))
			return
		}

		c.JSON(http.StatusCreated, response.SuccessResponse(newUser, "Successfully created a new user"))
	}
}

func (uc *UserController) GetUserByID() gin.HandlerFunc {
	return func (c *gin.Context) {
		id := c.Param("id")
		user, err := userColl.GetUserByID(id); if err != nil {
			c.JSON(http.StatusNotFound, response.ErrorResponse(err, "User not found"))
			return
		}

		c.JSON(http.StatusOK, response.SuccessResponse(user, "Successfully retrieved user by ID"))
	}
}

func (uc *UserController) UpdateUser() gin.HandlerFunc {
	return func (c *gin.Context) {
		id := c.Param("id")
		var user database.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(err, "Failed to update user"))
			return
		}

		updatedUser, err := userColl.UpdateUser(id, user); if err != nil {
			c.JSON(http.StatusNotFound, response.ErrorResponse(err, "User not found"))
			return
		}

		c.JSON(http.StatusOK, response.SuccessResponse(updatedUser, "Successfully updated user"))
	}
}

func (uc *UserController) DeleteUser() gin.HandlerFunc {
	return func (c *gin.Context) {
		id := c.Param("id")
		user, err := userColl.DeleteUser(id); if err != nil {
			c.JSON(http.StatusNotFound, response.ErrorResponse(err, "User not found"))
			return
		}

		c.JSON(http.StatusOK, response.SuccessResponse(user, "Successfully deleted user"))
	}
}

func (uc *UserController) LoginUser() gin.HandlerFunc {
	return func (c *gin.Context) {
		var user database.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(err, "Failed to login user"))
			return
		}

		user, err := userColl.LoginUser(user.Email, user.Password); if err != nil {
			c.JSON(http.StatusNotFound, response.ErrorResponse(nil, fmt.Sprintf("%v", err)))
			return
		}

		token := getToken(user)
		if token == "" {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse(nil, "Failed to generate token"))
			return
		}

		resp := Login{
			User: user,
			Token: token,
		}

		c.JSON(http.StatusOK, response.SuccessResponse(resp, "Successfully logged in"))
	}
}

func (uc *UserController) AuthMiddleware() gin.HandlerFunc {
	return func (c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse(nil, "Unauthorized"))
			c.Abort()
			return
		}

		token := strings.Split(authHeader, " ")[1]
		if token == "" {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse(nil, "Unauthorized"))
			c.Abort()
			return
		}

		email, err := verifyToken(token); if err != nil {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse(nil, "Unauthorized"))
			c.Abort()
			return
		}

		user, err := userColl.GetUserByEmail(email); if err != nil {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse(err, "Unauthorized"))
			return
		}
		c.Set("user", user)
		c.Next()
	}
}

func (uc *UserController) AuthUser() gin.HandlerFunc {
	return func (c *gin.Context) {
		user, ok := c.Get("user"); if !ok {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse(nil, "Unauthorized"))
			return
		}

		c.JSON(http.StatusOK, response.SuccessResponse(user, "Successfully authenticated user"))
	}
}

func (uc *UserController) SuggestWithGPT() gin.HandlerFunc {
	return func (c *gin.Context) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		fmt.Println("Memory usage before GPT call: ", m.Alloc/1024, "kb")

		_, ok := c.Get("user"); if !ok {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse(nil, "Unauthorized"))
			return
		}

		var body SuggestRequest
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(err, "Failed to parse input"))
			return
		}

		message, _ := c.GetQuery("message")
		httpClient := internalhttp.New()

		fmt.Println(utils.ComputeGPTPrompt(message))
		url := "https://api.openai.com/v1/chat/completions"
		headers := map[string]string{
			"Content-Type": "application/json",
			"Authorization": "Bearer " + config.GetEnv("OPEN_AI_API_KEY"),
		}

		data := utils.ComputeGPTPrompt(message)
		resp, err := httpClient.Post(url, headers, data); if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, response.ErrorResponse(err, "Something went wrong"))
			return
		}

		var r OpenAIResponse
		err = json.Unmarshal(resp, &r); if err != nil {
			fmt.Println("Failed to unmarshal response:", err)
			c.JSON(http.StatusInternalServerError, response.ErrorResponse(err, "Failed to parse response"))
			return
		}

		var m2 runtime.MemStats
		runtime.ReadMemStats(&m2)
		fmt.Println("Memory usage after GPT call: ", m2.Alloc/1024, "kb")

		c.JSON(http.StatusOK, response.SuccessResponse(r.Choices[0].Message.Content, "Successfully"))
	}
}

func (uc *UserController) TextToSpeech() gin.HandlerFunc {
	return func (c *gin.Context) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		fmt.Println("Memory usage before GPT call: ", m.Alloc/1024, "kb")

		_, ok := c.Get("user"); if !ok {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse(nil, "Unauthorized"))
			return
		}

		var body SuggestRequest
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(err, "Failed to parse input"))
			return
		}

		message := body.Message
		httpClient := internalhttp.New()

		fmt.Println(utils.ComputeGPTPrompt(message))
		url := "https://texttospeech.googleapis.com/v1/text:synthesize?key="+config.GetEnv("GOOGLE_API_KEY")
		headers := map[string]string{
			"Content-Type": "application/json",
		}

		data := utils.ComputeTTSInput(message)
		resp, err := httpClient.Post(url, headers, data); if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, response.ErrorResponse(err, "Something went wrong"))
			return
		}

		var r TTSResponse
		err = json.Unmarshal(resp, &r); if err != nil {
			fmt.Println("Failed to unmarshal response:", err)
			c.JSON(http.StatusInternalServerError, response.ErrorResponse(err, "Failed to parse response"))
			return
		}

		base64String := r.AudioContent

		go utils.SaveBase64ToDisk(base64String, "hello.mp3")

		c.JSON(http.StatusOK, response.SuccessResponse(
			map[string]interface{}{
				"message": "audio generated successfully",
			}, "Successfully"))
	}
}
