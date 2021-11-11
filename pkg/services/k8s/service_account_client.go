package k8s

import (
	"context"

	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// GetServiceAccount fetches a service account with the provided name in the namespace
func (kc KubeClient) GetServiceAccount(
	ctx context.Context,
	namespacedName types.NamespacedName,
) (coreV1.ServiceAccount, error) {
	kc.logger.Infof("fetching service account %v", namespacedName)

	serviceAccount := coreV1.ServiceAccount{}

	err := kc.client.Get(ctx, namespacedName, &serviceAccount)

	return serviceAccount, err
}

// DeleteServiceAccount deletes a service account in the cluster
func (kc KubeClient) DeleteServiceAccount(ctx context.Context, sa *coreV1.ServiceAccount) error {
	kc.logger.Infof("deleting ServiceAccount: %s/%s", sa.Namespace, sa.Name)

	if err := kc.client.Delete(ctx, sa); err != nil {
		kc.logger.Errorf("failed to delete ServiceAccount: %s/%s with error %v", sa.Namespace, sa.Name, err)
		return err
	}

	return nil
}

// CreateOrUpdateServiceAccount creates or updates a service account in the cluster
func (kc KubeClient) CreateOrUpdateServiceAccount(
	ctx context.Context,
	sa *coreV1.ServiceAccount,
	mutator controllerutil.MutateFn,
) (controllerutil.OperationResult, error) {
	kc.logger.Infof("creating or updating ServiceAccount: %s/%s", sa.Namespace, sa.Name)

	return ctrl.CreateOrUpdate(ctx, kc.client, sa, mutator)
}
