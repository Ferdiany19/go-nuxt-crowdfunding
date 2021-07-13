package handler

import (
	"net/http"
	"startup/helper"
	"startup/user"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	//  tangkap input dari user
	var input user.RegisterUserInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)

		// menampilkan pesan error menggunakan array errors
		errMsg := gin.H{"errors": errors}
		response := helper.ApiResponse("Failed to register account", http.StatusUnprocessableEntity, "error", errMsg)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// masukkan ke db
	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		response := helper.ApiResponse("Failed to register account", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// format response json
	formatter := user.FormatUser(newUser, "token")
	response := helper.ApiResponse("Account has been registered", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context) {
	var input user.LoginInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)

		// menampilkan pesan error menggunakan array errors
		errMsg := gin.H{"errors": errors}
		response := helper.ApiResponse("Login Failed!", http.StatusUnprocessableEntity, "error", errMsg)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedinUser, err := h.userService.Login(input)
	if err != nil {
		errMsg := gin.H{"errors": err.Error()}
		response := helper.ApiResponse("Login Failed!", http.StatusUnprocessableEntity, "error", errMsg)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// mengembalikan ke json
	formatter := user.FormatUser(loggedinUser, "token")
	response := helper.ApiResponse("Login Successfully!", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) CheckEmailAvailibility(c *gin.Context) {
	// ada input email dari user
	// input email di mapping ke struct input
	// struct input passing ke service
	// service akan manggil repository - email sudah ada atau belum?
	// repository - db 

	var input user.CheckEmailInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)

		// menampilkan pesan error menggunakan array errors
		errMsg := gin.H{"errors": errors}
		response := helper.ApiResponse("Email not found!", http.StatusUnprocessableEntity, "error", errMsg)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	
	IsEmailAvailable, err := h.service.IsEmailAvailable(input)
	if err != nil {
		errMsg := gin.H{"errors": "Server Error"}
		response := helper.ApiResponse("Email not found!", http.StatusUnprocessableEntity, "error", errMsg)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data := gin.H {
		"is_available": IsEmailAvailable,
	}

	var metaMessage := "Email has been registered"
	if IsEmailAvailable {
		metaMessage = "Email is avaiable"
	} 

	response := helper.ApiResponse(metaMessage, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}