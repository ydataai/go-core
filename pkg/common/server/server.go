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
	logger        logging.Logger
	Router        *gin.Engine
	httpServer    *http.Server
	configuration HTTPServerConfiguration
	readyzFunc    func() bool
}

// NewServer initializes a server
func NewServer(logger logging.Logger, configuration HTTPServerConfiguration) *Server {
	router := gin.Default()

	s := &Server{
		logger:        logger,
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
		s.logger.Infof("Server Running on [%v:%v]", s.configuration.Host, s.configuration.Port)
		if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			s.logger.Errorf("unexpected error while running server %v", err)
		}

		c := make(chan os.Signal, 3)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c

		s.logger.Infof("Shutdown Server ...")

		if err := s.httpServer.Shutdown(ctx); err != nil {
			s.logger.Fatal("Server forced to shutdown: ", err)
		}

		s.logger.Infof("Server exiting")
	}()
}

// RunSecurely when called starts the https server
func (s *Server) RunSecurely(ctx context.Context) {
	go func() {
		s.logger.Infof("Server Running on [%v:%v]", s.configuration.Host, s.configuration.Port)
		if err := s.httpServer.ListenAndServeTLS(s.configuration.CertificateFile, s.configuration.CertificateKeyFile); err != http.ErrServerClosed {
			s.logger.Errorf("unexpected error while running server %v", err.Error())
		}

		c := make(chan os.Signal, 2)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c

		s.logger.Infof("Shutdown Server ...")

		if err := s.httpServer.Shutdown(ctx); err != nil {
			s.logger.Fatalf("Server forced to shutdown: %v", err)
		}

		s.logger.Infof("Server exiting")
	}()
}

// AddHealthz creates a route to LivenessProbe
func (s *Server) AddHealthz(urls ...string) {
	url := firstStringOfArrayWithFallback(urls, s.configuration.HealthzEndpoint)

	s.Router.GET(url, s.healthz())
}

// AddReadyz creates a route to ReadinessProbe
func (s *Server) AddReadyz(status func() bool, urls ...string) {
	// Readyz probe is negative by default
	s.readyzFunc = status

	url := firstStringOfArrayWithFallback(urls, s.configuration.ReadyzEndpoint)

	s.Router.GET(url, s.readyz())
}
