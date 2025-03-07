// Package kubernetes is an util library to deal with kubernetes.
package kubernetes

import (
	"context"
	"time"

	"github.com/ydataai/go-core/pkg/common/logging"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	ctrlClient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
)

// Client is an interface that extends "ctrlClient.Client" and adds additional methods for starting and
// retrieving uncached objects.
type Client interface {
	ctrlClient.Client

	Start(ctx context.Context) error

	GetUncached(ctx context.Context, key ctrlClient.ObjectKey, obj ctrlClient.Object, opts ...ctrlClient.GetOption) error
	ListUncached(ctx context.Context, list ctrlClient.ObjectList, opts ...ctrlClient.ListOption) error
}

// Options are the possible options that can be configured in a kubernetes Client.
type Options struct {
	// Scheme is the scheme used to resolve runtime.Objects to GroupVersionKinds / Resources
	// Defaults to the kubernetes/client-go scheme.Scheme, but it's almost always better
	// idea to pass your own scheme in.  See the documentation in pkg/scheme for more information.
	Scheme *runtime.Scheme

	// SyncPeriod determines the minimum frequency at which watched resources are
	// reconciled. A lower period will correct entropy more quickly, but reduce
	// responsiveness to change if there are many watched resources. Change this
	// value only if you know what you are doing. Defaults to 10 hours if unset.
	// there will a 10 percent jitter between the SyncPeriod of all controllers
	// so that all controllers will not send list requests simultaneously.
	SyncPeriod *time.Duration

	// DisableCacheFor tells the client that, if any cache is used, to bypass it
	// for the given objects.
	DisableCacheFor []ctrlClient.Object
}

type client struct {
	ctrlClient.Client
	cluster   cluster.Cluster
	apiReader ctrlClient.Reader
	logger    logging.Logger
}

// NewClient instantiates a new kubernetes client with a logger and provided options
// This creates a cached client and an uncached reader.
func NewClient(config *rest.Config, logger logging.Logger, options Options) (Client, error) {
	httpClient, err := rest.HTTPClientFor(config)
	if err != nil {
		logger.Error(err, "Failed to create HTTP Client for config", config)
	}
	// Create the mapper provider
	mapper, err := apiutil.NewDynamicRESTMapper(config, httpClient)
	if err != nil {
		logger.Error(err, "Failed to get API Group-Resources")
		return nil, err
	}

	clientOptions := ctrlClient.Options{Scheme: options.Scheme, Mapper: mapper}

	apiReader, err := ctrlClient.New(config, clientOptions)
	if err != nil {
		return nil, err
	}

	cluster, err := cluster.New(config, func(o *cluster.Options) {
		o.Cache = cache.Options{
			Scheme:     options.Scheme,
			Mapper:     mapper,
			SyncPeriod: options.SyncPeriod,
		}
		o.Client = ctrlClient.Options{
			HTTPClient: httpClient,
			Scheme:     options.Scheme,
			Mapper:     mapper,
			Cache: &ctrlClient.CacheOptions{
				DisableFor: options.DisableCacheFor,
			},
		}
		o.HTTPClient = httpClient
	})
	if err != nil {
		return nil, err
	}

	return client{cluster.GetClient(), cluster, apiReader, logger}, nil
}

func (c client) Start(ctx context.Context) error {
	return c.cluster.Start(ctx)
}

func (c client) GetUncached(
	ctx context.Context, key ctrlClient.ObjectKey, obj ctrlClient.Object, opts ...ctrlClient.GetOption,
) error {
	return c.apiReader.Get(ctx, key, obj, opts...)
}

func (c client) ListUncached(ctx context.Context, list ctrlClient.ObjectList, opts ...ctrlClient.ListOption) error {
	return c.apiReader.List(ctx, list, opts...)
}
