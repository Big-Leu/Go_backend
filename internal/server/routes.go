package server

import (
	"kubequntumblock/controllers"
	"kubequntumblock/internal/auth"
	"kubequntumblock/middleware"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080","http://localhost:3000/"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE","PATCH"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))
	r.GET("/", s.HelloWorldHandler)
	r.GET("/health", s.healthHandler)
	r.POST("/adduser", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.GET("/getPods",middleware.RequireAuth,controllers.GetPods)
	r.POST("/createPod",middleware.RequireAuth,controllers.CreatePods)
	r.POST("/addEndpoint",middleware.RequireAuth,controllers.ExecCommandInPod)
	r.PATCH("/patchpod",middleware.RequireAuth,controllers.PatchPod)
	r.DELETE("/deletepod",middleware.RequireAuth,controllers.DeletePod)
	r.PATCH("/getLogs",middleware.RequireAuth,controllers.GetPodsLogs)
	r.GET("/auth/google",auth.GoogleAuth)
	r.GET("/auth/callback",auth.GoogleAuthCallbackFunction)
	r.GET("/auth/github",auth.GitHubAuth)
	r.GET("/github/auth/callback",auth.GitHubAuthCallbackFunction)
	r.GET("/logout",auth.Logout)
	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
