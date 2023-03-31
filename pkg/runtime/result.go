package runtime

import (
	"github.com/gin-gonic/gin"
	"github.com/ydataai/go-core/pkg/common/errors"
)

// Result represents a runtime type to be used across layers to return data and a status
// Is like an http Response, but it doesn't requires body decode
type Result struct {
	Status int
	Data   any
}

// NewResult is an helper method to build a new generic result
func NewResult(status int, data any) Result {
	return Result{
		Status: status,
		Data:   data,
	}
}

// ContinueResponse incapsulates some
func ContinueResponse(data any) Result {
	return NewResult(100, data)
}

func SuccessResponse(data any) Result {
	return NewResult(200, data)
}

func ErrorResponse(err errors.FabricError) Result {
	return NewResult(err.HTTPCode, err)
}

func (r Result) WriteTo(ctx *gin.Context) {
	if r.Data != nil {
		ctx.JSON(r.Status, r.Data)
		return
	}

	ctx.Status(r.Status)
}
