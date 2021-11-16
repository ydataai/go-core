package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/ydataai/go-core/pkg/common/logging"
	"github.com/ydataai/go-core/pkg/services/cloud"

	"github.com/gin-gonic/gin"
)

// CloudAdapterRESTController defines rest controller
type CloudAdapterRESTController struct {
	logger  logging.Logger
	service cloud.MeteringService
}

// NewCloudAdapterRESTController initializes rest controller
func NewCloudAdapterRESTController(
	logger logging.Logger,
	service cloud.MeteringService,
) CloudAdapterRESTController {
	return CloudAdapterRESTController{
		service: service,
		logger:  logger,
	}
}

// Boot ...
func (r CloudAdapterRESTController) Boot(router *gin.Engine) {
	router.GET("/healthz", r.healthCheck())
	router.POST("/metering/usageevents", r.createUsageEvent())
}

func (r CloudAdapterRESTController) healthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	}
}

func (r CloudAdapterRESTController) createUsageEvent() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, time.Second*10)
		defer cancel()

		ueb := cloud.UsageEventBatchReq{}
		if err := c.BindJSON(&ueb); err != nil {
			r.logger.Errorf("Error to parse UsageEventBatchReq request: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		res, err := r.service.CreateUsageEventBatch(ctx, ueb)
		if err != nil {
			r.logger.Errorf("Error to perform CreateUsageEventBatch: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	}
}
