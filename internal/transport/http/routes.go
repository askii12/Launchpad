package http

import (
	"github.com/askii12/launchpad/internal/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {

	userRepo := repository.NewUserRepository(db)
	userHandler := NewUserHandler(userRepo)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.POST("/users", userHandler.CreateUser)
	r.GET("/users", userHandler.GetUsers)
}
