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

type Client interface {
	ctrlClient.Client

	Start(ctx context.Context) error

	GetUncached(ctx context.Context, key ctrlClient.ObjectKey, obj ctrlClient.Object, opts ...ctrlClient.GetOption) error
	ListUncached(ctx context.Context, list ctrlClient.ObjectList, opts ...ctrlClient.ListOption) error
}

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
	cache     cache.Cache
	apiReader ctrlClient.Reader
	logger    logging.Logger
}

func NewClient(config *rest.Config, logger logging.Logger, options Options) (Client, error) {
	// Create the mapper provider
	mapper, err := apiutil.NewDynamicRESTMapper(config)
	if err != nil {
		logger.Error(err, "Failed to get API Group-Resources")
		return nil, err
	}

	// Create the cache for the cached read client and registering informers
	cacheOptions := cache.Options{
		Scheme: options.Scheme,
		Mapper: mapper,
		Resync: options.SyncPeriod,
	}

	cache, err := cache.New(config, cacheOptions)
	if err != nil {
		return nil, err
	}

	clientOptions := ctrlClient.Options{Scheme: options.Scheme, Mapper: mapper}

	apiReader, err := ctrlClient.New(config, clientOptions)
	if err != nil {
		return nil, err
	}

	clusterClient, err := cluster.DefaultNewClient(cache, config, clientOptions, options.DisableCacheFor...)
	if err != nil {
		return nil, err
	}

	return client{clusterClient, cache, apiReader, logger}, nil
}

func (c client) Start(ctx context.Context) error {
	return c.cache.Start(ctx)
}

func (c client) GetUncached(
	ctx context.Context, key ctrlClient.ObjectKey, obj ctrlClient.Object, opts ...ctrlClient.GetOption,
) error {
	return c.apiReader.Get(ctx, key, obj, opts...)
}

func (c client) ListUncached(ctx context.Context, list ctrlClient.ObjectList, opts ...ctrlClient.ListOption) error {
	return c.apiReader.List(ctx, list, opts...)
}
