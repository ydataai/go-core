package marketplace

import (
	"context"
)

// MeteringService represents the metering operations and it's responsible to handle
// each specific cloud provider business rules, data conversions and validations.
type MeteringService interface {
	// CreateUsageEvent creates a single usage event at the cloud provider metering service.
	CreateUsageEvent(ctx context.Context, req UsageEvent) (UsageEventResponse, error)
	// CreateUsageEventBatch creates a event batch  at the cloud provider metering service.
	CreateUsageEventBatch(ctx context.Context, req UsageEventBatch) (UsageEventBatchResponse, error)
}
