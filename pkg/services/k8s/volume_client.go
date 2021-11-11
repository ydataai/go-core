package k8s

import (
	"context"

	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// >>>>>>>>>>  PERSISTENT VOLUME <<<<<<<<<<

// GetPersistentVolume fetches a persistent volume in the cluster
func (kc KubeClient) GetPersistentVolume(
	ctx context.Context,
	namespacedName types.NamespacedName,
) (coreV1.PersistentVolume, error) {
	kc.logger.Infof("fetching persistent volume %v", namespacedName)
	pv := coreV1.PersistentVolume{}

	err := kc.client.Get(ctx, namespacedName, &pv)

	return pv, err
}

// DeletePersistentVolume attempts to delete a given PV
func (kc KubeClient) DeletePersistentVolume(ctx context.Context, pv *coreV1.PersistentVolume) error {
	kc.logger.Infof("deleting persistent volume: %v/%v", pv.Namespace, pv.Name)
	return kc.client.Delete(ctx, pv)
}

// CreateOrUpdatePersistentVolume attempts to create, if not existent, or to update an existent PersistentVolumeClaim
func (kc KubeClient) CreateOrUpdatePersistentVolume(
	ctx context.Context,
	pv *coreV1.PersistentVolume,
	mutator controllerutil.MutateFn,
) (controllerutil.OperationResult, error) {
	kc.logger.Infof("creating or updating persistent volume: %v/%v", pv.Namespace, pv.Name)
	return ctrl.CreateOrUpdate(ctx, kc.client, pv, mutator)
}

// >>>>>>>>>>  PERSISTENT VOLUME CLAIM <<<<<<<<<<

// GetPersistentVolumeClaim fetches a persistent volume claim in the cluster
func (kc KubeClient) GetPersistentVolumeClaim(
	ctx context.Context,
	namespacedName types.NamespacedName,
) (coreV1.PersistentVolumeClaim, error) {
	kc.logger.Infof("fetching persistent volume claim: %v", namespacedName)
	claim := coreV1.PersistentVolumeClaim{}

	err := kc.client.Get(ctx, namespacedName, &claim)

	return claim, err
}

// DeletePersistentVolumeClaim attempts to delete a given PVC
func (kc KubeClient) DeletePersistentVolumeClaim(ctx context.Context, pvc *coreV1.PersistentVolumeClaim) error {
	kc.logger.Infof("deleting persistent volume: %v/%v", pvc.Namespace, pvc.Name)
	return kc.client.Delete(ctx, pvc)
}

// CreateOrUpdatePersistentVolumeClaim ...
func (kc KubeClient) CreateOrUpdatePersistentVolumeClaim(
	ctx context.Context,
	pvc *coreV1.PersistentVolumeClaim,
	mutator controllerutil.MutateFn,
) (controllerutil.OperationResult, error) {
	kc.logger.Infof("creating or updating persistent volume claim: %v/%v", pvc.Namespace, pvc.Name)
	return ctrl.CreateOrUpdate(ctx, kc.client, pvc, mutator)
}
