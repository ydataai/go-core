// Package metering provides client and tools to interact with the adapter server
package metering

import (
	"context"
	"fmt"
	"net/http"

	coreHTTP "github.com/ydataai/go-core/pkg/http"
)

// MeteringClient represents the metering operations and it's responsible to handle
// each specific cloud provider business rules, data conversions and validations.
type MeteringClient interface {
	// CreateUsageEvent creates a single usage event at the cloud provider metering service.
	CreateUsageEvent(ctx context.Context, req UsageEvent) (UsageEventResponse, error)
	// CreateUsageEventBatch creates a event batch  at the cloud provider metering service.
	CreateUsageEventBatch(ctx context.Context, req UsageEventBatch) (UsageEventBatchResponse, error)
}

type MeteringClientOptions struct {
	BaseURL string
}

// Adapter usually runs on the same machine as a side car on port 8081
const defaultBaseURL = "http://localhost:8081"

const (
	usageEvent      = "usageEvent"
	batchUsageEvent = "batchUsageEvent"
)

type meteringClient struct {
	pl      coreHTTP.Pipeline
	options MeteringClientOptions
}

func NewMeteringClient(options *MeteringClientOptions) MeteringClient {
	if options == nil {
		defaultOptions := defaultOptions()
		options = &defaultOptions
	}

	// Work in progress
	pl := coreHTTP.NewPipeline()

	return meteringClient{
		pl:      pl,
		options: *options,
	}
}

func (c meteringClient) CreateUsageEvent(ctx context.Context, req UsageEvent) (UsageEventResponse, error) {
	result, err := sendRequest[UsageEvent, UsageEventResponse](ctx, c.pl, c.options.BaseURL, usageEvent, req)
	return *result, err
}

func (c meteringClient) CreateUsageEventBatch(
	ctx context.Context, req UsageEventBatch,
) (UsageEventBatchResponse, error) {
	result, err := sendRequest[UsageEventBatch, UsageEventBatchResponse](
		ctx, c.pl, c.options.BaseURL, batchUsageEvent, req)
	return *result, err
}

func sendRequest[T, V any](
	ctx context.Context, pl coreHTTP.Pipeline, baseURL string, path string, obj T,
) (*V, error) {
	endpoint := coreHTTP.JoinPaths(baseURL, "metering", path)
	request, err := coreHTTP.NewRequest(ctx, http.MethodPost, endpoint)
	if err != nil {
		return nil, err
	}

	if err := request.EncodeAsJSON(obj); err != nil {
		return nil, err
	}

	resp, err := pl.Do(request)
	if err != nil {
		return nil, err
	}

	if !resp.HasStatusCode(http.StatusAccepted, http.StatusOK) {
		return nil, invalidStatusCodeError(resp.Response)
	}

	result := new(V)
	if err := resp.DecodeJSON(result); err != nil {
		return nil, err
	}

	return result, nil
}

func defaultOptions() MeteringClientOptions {
	return MeteringClientOptions{
		BaseURL: defaultBaseURL,
	}
}

func invalidStatusCodeError(response *http.Response) error {
	return fmt.Errorf("request failed with error %s", response.Status)
}
