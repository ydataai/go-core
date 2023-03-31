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

// Server represents a HTTP Server with gin router and health probes
type Server interface {
	Router() *gin.Engine

	AddHealthz(urls ...string)
	AddReadyz(status *func() bool, urls ...string)

	Run(ctx context.Context, readyCallbacks ...func())
	RunSecurely(ctx context.Context, readyCallbacks ...func())

	// NamespaceValidation checks if namespace is in the query params
	NamespaceValidation() gin.HandlerFunc
}

var defaultReadyzFunc = func() bool { return true }

// Server defines a struct for server
type server struct {
	configuration HTTPServerConfiguration
	httpServer    *http.Server
	logger        logging.Logger

	router *gin.Engine

	readyzFunc func() bool
}

// NewServer initializes a server
func NewServer(logger logging.Logger, configuration HTTPServerConfiguration) Server {
	router := gin.Default()

	s := &server{
		configuration: configuration,
		logger:        logger,
		router:        router,
		readyzFunc:    defaultReadyzFunc,
	}

	s.router.Use(
		s.tracing(),
		s.setUserID(),
	)

	return s
}

func (s *server) Router() *gin.Engine {
	return s.router
}

// Run when called starts the server
// warning: once the Run is called, you cannot modify the Handle in http.Server.
func (s *server) Run(ctx context.Context, readyCallbacks ...func()) {
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
func (s *server) RunSecurely(ctx context.Context, readyCallbacks ...func()) {
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
func (s *server) AddHealthz(urls ...string) {
	url := firstStringOfArrayWithFallback(urls, s.configuration.HealthzEndpoint)

	s.router.GET(url, s.healthz())
}

// AddReadyz creates a route to ReadinessProbe
func (s *server) AddReadyz(status *func() bool, urls ...string) {
	// Readyz probe is negative by default
	if status != nil {
		s.readyzFunc = *status
	}

	url := firstStringOfArrayWithFallback(urls, s.configuration.ReadyzEndpoint)

	s.router.GET(url, s.readyz())
}

func (s *server) httpServerSetup() {
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.configuration.Host, s.configuration.Port),
		Handler: s.router,
	}
}

// waitingToBeReady executes a callback when the server is ready.
func (s *server) waitingToBeReady(readyCallbacks ...func()) {
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

func (s *server) shutdown(ctx context.Context) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	s.logger.Infof("Shutdown Server ...")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Fatalf("Server forced to shutdown: %v", err)
	}

	s.logger.Infof("Server exiting")
}
