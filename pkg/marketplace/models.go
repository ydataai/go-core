package marketplace

import "time"

// UsageEventReq represents an usage event for metering purpose
type UsageEvent struct {
	DimensionID string    `json:"dimensionId"`
	Quantity    float32   `json:"quantity"`
	StartAt     time.Time `json:"startAt"`
}

// UsageEventRes represents an usage event response
type UsageEventResponse struct {
	UsageEventID string `json:"usageEventId"`
	DimensionID  string `json:"dimensionId"`
	Status       string `json:"status"`
}

// UsageEventBatchReq a type to represent the usage metering batch events request
type UsageEventBatch struct {
	Events []UsageEvent `json:"events"`
}

// UsageEventBatchRes a type to represent the usage metering event batch response
type UsageEventBatchResponse struct {
	Result []UsageEventResponse `json:"result"`
}
