package k8s

import (
	"context"

	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// GetService fetches a service in the cluster
func (kc KubeClient) GetService(ctx context.Context, namespacedName types.NamespacedName) (coreV1.Service, error) {
	kc.logger.Infof("fetching Service: %v", namespacedName)
	service := coreV1.Service{}

	err := kc.client.Get(ctx, namespacedName, &service)

	return service, err
}

// DeleteService deletes a service in the cluster
func (kc KubeClient) DeleteService(ctx context.Context, svc *coreV1.Service) error {
	kc.logger.Infof("deleting Service: %s/%s", svc.Namespace, svc.Name)

	if err := kc.client.Delete(ctx, svc); err != nil {
		kc.logger.Errorf("while deleting Service: %s/%s", svc.Namespace, svc.Name)
		return err
	}

	return nil
}

// CreateOrUpdateService creates a service in the cluster
func (kc KubeClient) CreateOrUpdateService(
	ctx context.Context,
	service *coreV1.Service,
	mutator controllerutil.MutateFn,
) (controllerutil.OperationResult, error) {
	kc.logger.Infof("creating or updating Service: %s/%s", service.Namespace, service.Name)

	return ctrl.CreateOrUpdate(ctx, kc.client, service, mutator)
}
