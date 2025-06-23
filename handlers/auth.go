package handlers

import (
	"net/http"
	"job-api/config"
	"job-api/models"
	"job-api/utils"
	"github.com/gin-gonic/gin"
)

type SignupRequest struct {
	Name     string           `json:"name" validate:"required,alpha"`
	Email    string           `json:"email" validate:"required,email"`
	Password string           `json:"password" validate:"required,min=8,containsuppercase,containslowercase,containsdigit,containsspecial"`
	Role     models.UserRole  `json:"role" validate:"required,oneof=applicant company"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResponse{
			Success: false,
			Message: "Invalid request data",
			Object:  nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResponse{
			Success: false,
			Message: "Validation failed",
			Object:  nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := config.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, models.BaseResponse{
			Success: false,
			Message: "User already exists",
			Object:  nil,
			Errors:  []string{"Email already registered"},
		})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponse{
			Success: false,
			Message: "Failed to process password",
			Object:  nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	// Create user
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponse{
			Success: false,
			Message: "Failed to create user",
			Object:  nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	// Remove password from response
	user.Password = ""

	c.JSON(http.StatusCreated, models.BaseResponse{
		Success: true,
		Message: "User created successfully",
		Object:  user,
	})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResponse{
			Success: false,
			Message: "Invalid request data",
			Object:  nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResponse{
			Success: false,
			Message: "Validation failed",
			Object:  nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	// Find user
	var user models.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, models.BaseResponse{
			Success: false,
			Message: "User not found",
			Object:  nil,
			Errors:  []string{"Invalid credentials"},
		})
		return
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, models.BaseResponse{
			Success: false,
			Message: "Incorrect password",
			Object:  nil,
			Errors:  []string{"Invalid credentials"},
		})
		return
	}

	// Generate JWT
	token, err := utils.GenerateJWT(user.ID, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponse{
			Success: false,
			Message: "Failed to generate token",
			Object:  nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, models.BaseResponse{
		Success: true,
		Message: "Login successful",
		Object:  gin.H{"token": token, "user": gin.H{"id": user.ID, "name": user.Name, "email": user.Email, "role": user.Role}},
	})
}
