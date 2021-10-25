package grpc

import (
	context "context"

	"github.com/sirupsen/logrus"
	"github.com/ydataai/go-core/pkg/services/cloud"
)

// CloudAdaptergRPCController is used to implement cloud provider operations.
type CloudAdaptergRPCController struct {
	UnimplementedMeteringServiceServer
	service cloud.MeteringService
	logger  *logrus.Logger
}

// NewCloudAdaptergRPCController creates a new CloudAdaptergRPCController instance
func NewCloudAdaptergRPCController(service cloud.MeteringService, logger *logrus.Logger) CloudAdaptergRPCController {
	return CloudAdaptergRPCController{
		service: service,
		logger:  logger,
	}
}

// CreateUsageEvent creates a single usage event at the cloud provider metering service.
func (c CloudAdaptergRPCController) CreateUsageEvent(ctx context.Context, req *CreateUsageEventReq) (*CreateUsageEventRes, error) {
	c.logger.Infof("Received CreateUsageEventReq: %v", req)

	res, err := c.service.CreateUsageEvent(ctx, cloud.UsageEventReq{
		ResourceID: req.ResourceId,
		Quantity:   req.Quantity,
		StartAt:    req.StartAt.AsTime(),
	})
	if err != nil {
		c.logger.Errorf("Error to create cloud.UsageEventReq: %v", err)
		return nil, err
	}

	return &CreateUsageEventRes{
		ResourceId:   res.ResourceID,
		UsageEventId: res.UsageEventID,
		Status:       res.Status}, nil
}

// CreateUsageEventBatch creates a event batch  at the cloud provider metering service.
func (c CloudAdaptergRPCController) CreateUsageEventBatch(ctx context.Context, req *CreateUsageEventBatchReq) (*CreateUsageEventBatchRes, error) {
	c.logger.Infof("Received CreateUsageEventBatchReq: %v", req)

	events := make([]cloud.UsageEventReq, len(req.Request))
	for _, event := range req.Request {
		events = append(events, cloud.UsageEventReq{
			ResourceID: event.ResourceId,
			Quantity:   event.Quantity,
			StartAt:    event.StartAt.AsTime(),
		})
	}

	resp, err := c.service.CreateUsageEventBatch(ctx, cloud.UsageEventBatchReq{Request: events})
	if err != nil {
		c.logger.Errorf("Error to process cloud.UsageEventBatchReq: %v", err)
		return nil, err
	}

	results := make([]*CreateUsageEventRes, len(resp.Result))
	for _, result := range resp.Result {
		results = append(results, &CreateUsageEventRes{
			ResourceId:   result.ResourceID,
			UsageEventId: result.UsageEventID,
			Status:       result.Status,
		})
	}
	return &CreateUsageEventBatchRes{Result: results}, nil
}
