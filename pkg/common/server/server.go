package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
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

	s.Router = router

	return s
}

// Run when called starts the server
// warning: once the Run is called, you cannot modify the Handle in http.Server.
func (s *Server) Run(ctx context.Context, readyCallbacks ...func()) {
	s.httpServerSetup()

	go func() {
		s.logger.Infof("Server Running on [%v:%v]", s.configuration.Host, s.configuration.Port)
		if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			s.logger.Errorf("unexpected error while running server %v", err)
		}
	}()

	s.waitingToBeReady(readyCallbacks...)

	go s.shutdown(ctx)
}

// RunSecurely when called starts the https server
// warning: once the Run is called, you cannot modify the Handle in http.Server.
func (s *Server) RunSecurely(ctx context.Context, readyCallbacks ...func()) {
	s.httpServerSetup()

	go func() {
		s.logger.Infof("Server Running on [%v:%v]", s.configuration.Host, s.configuration.Port)
		if err := s.httpServer.ListenAndServeTLS(s.configuration.CertificateFile, s.configuration.CertificateKeyFile); err != http.ErrServerClosed {
			s.logger.Errorf("unexpected error while running server %v", err.Error())
		}
	}()

	s.waitingToBeReady(readyCallbacks...)

	go s.shutdown(ctx)
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

func (s *Server) httpServerSetup() {
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.configuration.Host, s.configuration.Port),
		Handler: s.Router,
	}
}

// waitingToBeReady executes a callback when the server is ready.
func (s *Server) waitingToBeReady(readyCallbacks ...func()) {
	for {
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(s.configuration.Host, strconv.Itoa(s.configuration.Port)), s.configuration.RequestTimeout)
		if err != nil {
			s.logger.Debug("Server is not ready yet!")
			continue
		}
		if conn != nil {
			s.logger.Infof("Server is ready!")
			break
		}
	}

	for _, callback := range readyCallbacks {
		callback()
	}
}

func (s *Server) shutdown(ctx context.Context) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	s.logger.Infof("Shutdown Server ...")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Fatalf("Server forced to shutdown: %v", err)
	}

	s.logger.Infof("Server exiting")
}
