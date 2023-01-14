package handlershttp

import (
	"context"
	"net/http"

	"github.com/labstack/echo"
)

type Server struct {
	Router *echo.Echo
	port   string
}

func NewServer(ctx context.Context, port string) *Server {
	return &Server{
		Router: echo.New(),
		port:   port,
	}
}

func (s *Server) Setup(ctx context.Context) error {
	s.Router.GET("/healthcheck", func(c echo.Context) error {
		type HealthResponse struct {
			Status string `json:"status"`
		}
		u := &HealthResponse{
			Status: "OK",
		}
		return c.JSON(http.StatusOK, u)
	})

	return nil
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.Router.Start(":" + s.port); err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.Router.Shutdown(ctx)
}
