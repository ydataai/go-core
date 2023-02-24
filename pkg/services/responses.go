package services

import (
	"github.com/gin-gonic/gin"
	"github.com/ydataai/go-core/pkg/common/errors"
)

type ServiceResponse struct {
	Status int
	Data   interface{}
}

func ContinueResponse(data interface{}) ServiceResponse {
	return ServiceResponse{
		Status: 100,
		Data:   data,
	}
}

func SuccessResponse(data interface{}) ServiceResponse {
	return ServiceResponse{
		Status: 200,
		Data:   data,
	}
}

func ErrorResponse(err errors.FabricError) ServiceResponse {
	return ServiceResponse{
		Status: err.HTTPCode,
		Data:   err,
	}
}

func (s ServiceResponse) WriteTo(ctx *gin.Context) {
	if s.Data != nil {
		ctx.JSON(s.Status, s.Data)
	}

	ctx.Status(s.Status)
}
