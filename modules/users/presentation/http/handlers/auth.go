package handlers

import (
	"net/http"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/configs"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/middleware"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/presentation/http/dto"
	"github.com/gin-gonic/gin"
)

func (h *UserHandler) Login(c *gin.Context) {
	var loginDto dto.LoginDto
	if err := c.ShouldBindJSON(&loginDto); err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"Invalid email or password format",
			err,
			http.StatusBadRequest,
		))
		return
	}

	// Look up user by email
	user, err := h.FindUserByEmailUsecase.Execute(c.Request.Context(), loginDto.Email)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email"})
		return
	}

	// Check password (assumes hashed password in DB)
	if !user.ComparePassword(loginDto.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Generate JWT
	jwtValidator := middleware.NewJWTValidator(configs.GetEnv("JWT_SECRET", ""))
	orgID := ""
	if user.OrganizationID != "" {
		orgID = user.OrganizationID
	}
	token, err := jwtValidator.GenerateToken(user.ID.Hex(), user.Role, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
