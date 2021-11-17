package cloud

import (
	"context"
	"time"
)

// UsageEventReq represents an usage event for metering purpose
type UsageEventReq struct {
	DimensionID string    `json:"dimensionId"`
	Quantity    float32   `json:"quantity"`
	StartAt     time.Time `json:"startAt"`
}

// UsageEventRes represents an usage event response
type UsageEventRes struct {
	UsageEventID string `json:"usageEventId"`
	DimensionID  string `json:"dimensionId"`
	Status       string `json:"status"`
}

// UsageEventBatchReq a type to represent the usage metering batch events request
type UsageEventBatchReq struct {
	Request []UsageEventReq `json:"request"`
}

// UsageEventBatchRes a type to represent the usage metering event batch response
type UsageEventBatchRes struct {
	Result []UsageEventRes `json:"result"`
}

// MeteringService represents the metering operations and it's responsible to handle
// each specific cloud provider business rules, data conversions and validations.
type MeteringService interface {
	// CreateUsageEvent creates a single usage event at the cloud provider metering service.
	CreateUsageEvent(ctx context.Context, req UsageEventReq) (UsageEventRes, error)
	// CreateUsageEventBatch creates a event batch  at the cloud provider metering service.
	CreateUsageEventBatch(ctx context.Context, req UsageEventBatchReq) (UsageEventBatchRes, error)
}
