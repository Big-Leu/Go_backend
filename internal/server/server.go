package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
    	"kubequntumblock/pkg/controllers"
	_ "github.com/joho/godotenv/autoload"

	"kubequntumblock/internal/database"
)

type Server struct {
	port int
 
	db database.Service

	K controllers.KubeService
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,

		db: database.New(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
