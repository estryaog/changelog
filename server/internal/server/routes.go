package server

import (
	"net/http"
	"time"

	"github.com/estryaog/changelog/internal/handler"
	"github.com/estryaog/changelog/internal/types"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()
	configureCors(r)

	r.POST("/register", s.Register)
	r.POST("/login", s.Login)
	r.GET("/me", s.Me)

	return r
}

func configureCors(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}

func (s *Server) Register(ctx *gin.Context) {
	var user *types.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json request"})
		return
	}

	user.Id = uuid.New().String()
	user.IsAdmin = false

	hashedPassword, err := handler.HashPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.Password = string(hashedPassword)

	emailExists, err := s.db.IsKeyValueExist("users", "email", user.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if emailExists {
		ctx.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	}

	createdUser, err := s.UserStore.CreateUser(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdUser)
}

func (s *Server) Login(ctx *gin.Context) {
	var user *types.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json request"})
		return
	}

	existingUser, err := s.UserStore.GetUser(ctx, "email", user.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "We couldn't find an account with that email address"})
		return
	}

	err = handler.CompareHashAndPassword(existingUser.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := handler.CreateJWT(existingUser.Id, existingUser.IsAdmin)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (s *Server) Me(ctx *gin.Context) {
	tokenString := ctx.GetHeader("Authorization")
	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
		return
	}

	claims, err := handler.ValidateJWT(tokenString)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	user, err := s.UserStore.GetUser(ctx, "id", claims.Id)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "We couldn't find an account with that id"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
