package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"bkoiki950/go-store/api/database"
	"bkoiki950/go-store/api/response"
)

type Login struct {
	User database.User `json:"user"`
	Token string `json:"token"`
}

func CreateUserController() *UserController {
	var userCtrl IUserController = &UserController{}
	return userCtrl.(*UserController)
}

var userColl = database.UserCollection()

func (uc *UserController) GetUsers(c *gin.Context) {
	users, err := userColl.GetAllUsers(database.User{}); if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err, "Failed to retrieve all users"))
		return
	}

	fmt.Println(users)

	c.JSON(http.StatusOK, response.SuccessResponse(users, "Successfully retrieved all users"))
}

func (uc *UserController) CreateUser(c *gin.Context) {
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

func (uc *UserController) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := userColl.GetUserByID(id); if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(err, "User not found"))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(user, "Successfully retrieved user by ID"))
}

func (uc *UserController) UpdateUser(c *gin.Context) {
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

func (uc *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	user, err := userColl.DeleteUser(id); if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(err, "User not found"))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(user, "Successfully deleted user"))
}

func (uc *UserController) LoginUser(c *gin.Context) {
	var user database.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err, "Failed to login user"))
		return
	}

	user, err := userColl.LoginUser(user.Email, user.Password); if err != nil {
		fmt.Println("Error: ", err)
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

