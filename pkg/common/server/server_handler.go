package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// healthz is a liveness probe
func (s *Server) healthz() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Status(http.StatusNoContent)
	}
}

// readyz is a readiness probe
func (s *Server) readyz() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !s.readyzFunc() {
			http.Error(ctx.Writer, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		ctx.Status(http.StatusNoContent)
	}
}
