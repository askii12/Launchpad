package http

import (
	"github.com/askii12/launchpad/internal/dto"
	"github.com/askii12/launchpad/internal/models"
	"github.com/askii12/launchpad/internal/repository"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) CreateUser(c *gin.Context) {

	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.repo.Create(&user); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, dto.UserResponse{
		ID:    user.ID,
		Email: user.Email,
	})
}
func (h *UserHandler) GetUsers(c *gin.Context) {

	users, err := h.repo.GetAll()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var res []dto.UserResponse
	for _, u := range users {
		res = append(res, dto.UserResponse{
			ID:    u.ID,
			Email: u.Email,
		})
	}

	c.JSON(200, res)
}
