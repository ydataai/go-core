package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/ydataai/go-core/pkg/common/logging"
)

// Server defines a struct for server
type Server struct {
	log           logging.Logger
	Router        *gin.Engine
	httpServer    *http.Server
	configuration HTTPServerConfiguration
}

// NewServer initializes a server
func NewServer(log logging.Logger, configuration HTTPServerConfiguration) *Server {
	router := gin.Default()

	s := &Server{
		log:           log,
		configuration: configuration,
	}

	router.Use(
		s.tracing(),
		s.setUserID(),
	)

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", configuration.Host, configuration.Port),
		Handler: router,
	}

	s.Router = router

	return s
}

// Run when called starts the server
func (s *Server) Run(ctx context.Context) {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			s.log.Errorf("unexpected error while running server %v", err)
		}

		c := make(chan os.Signal, 3)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c

		s.log.Infof("Shutdown Server ...")

		if err := s.httpServer.Shutdown(ctx); err != nil {
			s.log.Fatal("Server forced to shutdown: ", err)
		}

		s.log.Infof("Server exiting")
	}()
}

// RunSecurely when called starts the https server
func (s *Server) RunSecurely(ctx context.Context) {
	go func() {
		if err := s.httpServer.ListenAndServeTLS(s.configuration.CertificateFile, s.configuration.CertificateKeyFile); err != http.ErrServerClosed {
			s.log.Errorf("unexpected error while running server %v", err.Error())
		}

		c := make(chan os.Signal, 2)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c

		s.log.Infof("Shutdown Server ...")

		if err := s.httpServer.Shutdown(ctx); err != nil {
			s.log.Fatalf("Server forced to shutdown: %v", err)
		}

		s.log.Infof("Server exiting")
	}()
}

// UseHealthCheck creates a new HealthCheck route
func (s *Server) UseHealthCheck() {
	s.Router.GET(s.configuration.HealthCheckEndpoint, s.healthz())
}
