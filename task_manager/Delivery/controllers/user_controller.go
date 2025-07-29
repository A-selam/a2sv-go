package task_controllers

import (
	"net/http"
	domain "task_manager/Domain"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserUsecase domain.UserUsecase
}

func NewUserController(uu domain.UserUsecase) *UserController{
	return &UserController{
		UserUsecase: uu,
	}
}

// POST /register
func (uc *UserController) RegisterUser(c *gin.Context) {
	var user domain.User

	// Bind incoming JSON to user struct
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Call the service layer to register the user
	if err := uc.UserUsecase.Register(c, user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// POST /login
func (uc *UserController) LoginUser(c *gin.Context) {
	var user domain.Login

	// Bind incoming JSON to user struct
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Call the service layer to login and get token
	token, err := uc.UserUsecase.Login(c, user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
