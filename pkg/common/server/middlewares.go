package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s Server) tracing() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := uuid.New().String()
		if ctx.Request.Header.Get("X-Request-Id") != "" {
			requestID = ctx.Request.Header.Get("X-Request-Id")
		}

		s.logger.Infof("Path: %v", ctx.Request.URL.Path)
		s.logger.Infof("X-Request-Id: %v", requestID)
		ctx.Set("X-Request-Id", requestID)
		ctx.Next()
	}
}

// NamespaceValidation checks if namespace is in the query params
func (s Server) NamespaceValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Query("ns")

		if namespace == "" {
			errM := "missing namespace as query param."
			s.logger.Error(errM)
			c.JSON(http.StatusBadRequest, errM)
			c.Abort()
			return
		}

		c.Next()
	}
}

func (s Server) setUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userID", s.configuration.UserID)
		c.Next()
	}
}

// healthz is a liveness probe
func (s *Server) healthz() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Status(http.StatusNoContent)
	}
}

// readyz is a readiness probe
func (s *Server) readyz() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !s.readyzAvailable.Load().(bool) {
			http.Error(ctx.Writer, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		ctx.Status(http.StatusNoContent)
	}
}

// SetReadyzState enables/disables readyz
func (s *Server) SetReadyzState(state bool) {
	s.readyzAvailable.Store(state)
}
