package server

import (
	"net/http"
	"time"

	"github.com/estryaog/changelog/internal/types"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
	var reqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if reqBody.Email == "" || reqBody.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}

	user := &types.User{
		Id: uuid.New().String(),
		Email:    reqBody.Email,
		IsAdmin: false,
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.Password = string(hashedPassword)

	createdUser, err := s.UserStore.CreateUser(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdUser)
}

func (s *Server) Login(ctx *gin.Context) {
	var reqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if reqBody.Email == "" || reqBody.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}

	user := &types.User{
		Email:    reqBody.Email,
		Password: reqBody.Password,
	}

	existingUser, err := s.UserStore.GetUser(ctx, user.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	ctx.JSON(http.StatusOK, existingUser)
}

func (s *Server) Me(ctx *gin.Context) {
	/*
	 * ToDo: Implement logic to get the current user from the jwt token
	 */

	ctx.JSON(http.StatusOK, "implement logic")
}