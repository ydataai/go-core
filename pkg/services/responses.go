// Package service contains methods used by internal components
package services

import (
	"github.com/gin-gonic/gin"
	"github.com/ydataai/go-core/pkg/common/errors"
)

// ServiceResponse is a struct for all the responses made by the services
type ServiceResponse struct {
	Status int
	Data   interface{}
}

// ContinueResponse helper method to build a continue type ServiceResponse
func ContinueResponse(data interface{}) ServiceResponse {
	return ServiceResponse{
		Status: 100,
		Data:   data,
	}
}

// SuccessResponse helper method to build a success type ServiceResponse
func SuccessResponse(data interface{}) ServiceResponse {
	return ServiceResponse{
		Status: 200,
		Data:   data,
	}
}

// ErrorResponse helper method to build an error type ServiceResponse
func ErrorResponse(err errors.FabricError) ServiceResponse {
	return ServiceResponse{
		Status: err.HTTPCode,
		Data:   err,
	}
}

// WriteTo implements the gin interface to write a ServiceResponse into an gin request/response
func (s ServiceResponse) WriteTo(ctx *gin.Context) {
	if s.Data != nil {
		ctx.JSON(s.Status, s.Data)
	}

	ctx.Status(s.Status)
}
