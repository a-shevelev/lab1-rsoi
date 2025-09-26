package server

import (
	"fmt"
	handlers "lab1-rsoi/internal/handlers/http/v1"
	"lab1-rsoi/internal/repo"
	"lab1-rsoi/internal/service"
	"lab1-rsoi/pkg/postgres"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	Host      string `envconfig:"HOST" required:"true"`
	Port      int    `envconfig:"PORT" required:"true"`
	DB        postgres.Client
	GinRouter *gin.Engine
}

func New(dbc postgres.Client, host string, port int) (*Server, error) {
	s := &Server{
		Host:      host,
		Port:      port,
		DB:        dbc,
		GinRouter: gin.Default(),
	}

	if err := s.initRoutes(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Server) initRoutes() error {
	s.GinRouter.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "pong"})
	})

	if err := s.InitDocsRoutes(); err != nil {
		log.Info("Docs routes initialization failed")
	}

	v1 := s.GinRouter.Group("/api/v1")

	personRepo := repo.New(s.DB)
	personService := service.New(personRepo)
	personHandler := handlers.New(personService)
	personHandler.RegisterRoutes(v1)

	return nil
}

func (s *Server) Run() error {
	return s.GinRouter.Run(fmt.Sprintf("%s:%d", s.Host, s.Port))
}
