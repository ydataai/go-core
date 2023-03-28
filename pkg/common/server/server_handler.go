package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// healthz is a liveness probe
func (s *server) healthz() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Status(http.StatusNoContent)
	}
}

// readyz is a readiness probe
func (s *server) readyz() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !s.readyzFunc() {
			http.Error(ctx.Writer, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		ctx.Status(http.StatusNoContent)
	}
}

// firstStringOfArrayWithFallback is a handle to get the first position of the array with fallback
func firstStringOfArrayWithFallback(list []string, fallback string) string {
	if len(list) > 0 {
		return list[0]
	}
	return fallback
}
