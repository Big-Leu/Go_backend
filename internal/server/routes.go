package server

import (
	"kubequntumblock/controllers"
	"kubequntumblock/middleware"
	"net/http"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()
    r.Use(cors.Default())
	r.GET("/", s.HelloWorldHandler)
	r.GET("/health", s.healthHandler)
	r.POST("/adduser", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.GET("/getPods",controllers.GetPods)
	r.POST("/createPod",controllers.CreatePods)
	r.POST("/addEndpoint",controllers.ExecCommandInPod)
	r.PATCH("/patchpod",controllers.PatchPod)
	r.DELETE("/deletepod",controllers.DeletePod)
	r.PATCH("/getLogs",controllers.GetPodsLogs)
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
