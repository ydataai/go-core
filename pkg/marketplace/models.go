package marketplace

import "time"

// UsageEvent represents an usage event for metering purpose
type UsageEvent struct {
	DimensionID string    `json:"dimensionId"`
	Quantity    float32   `json:"quantity"`
	StartAt     time.Time `json:"startAt"`
}

// UsageEventResponse represents an usage event response
type UsageEventResponse struct {
	UsageEventID string `json:"usageEventId"`
	DimensionID  string `json:"dimensionId"`
	Status       string `json:"status"`
}

// UsageEventBatch a type to represent the usage metering batch events request
type UsageEventBatch struct {
	Events []UsageEvent `json:"events"`
}

// UsageEventBatchResponse a type to represent the usage metering event batch response
type UsageEventBatchResponse struct {
	Result []UsageEventResponse `json:"result"`
}
