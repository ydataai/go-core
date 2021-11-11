package k8s

import (
	"context"

	coreV1 "k8s.io/api/core/v1"
	apiErrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// GetPod fetches a pod on cluster
func (kc KubeClient) GetPod(ctx context.Context, namespacedName types.NamespacedName) (coreV1.Pod, error) {
	kc.logger.Infof("fetching Pod: %v", namespacedName)
	pod := coreV1.Pod{}

	err := kc.client.Get(ctx, namespacedName, &pod)

	return pod, err
}

// DeletePod deletes a pod on the cluster
func (kc KubeClient) DeletePod(ctx context.Context, pod *coreV1.Pod) error {
	kc.logger.Infof("deleting Pod: %s/%s", pod.Namespace, pod.Name)

	if err := kc.client.Delete(ctx, pod); err != nil && !apiErrs.IsNotFound(err) {
		kc.logger.Errorf("failed to delete Pod: %s/%s with error %v", pod.Namespace, pod.Name, err)
		return err
	}

	return nil
}

// CreateOrUpdatePod attempts to create, if not existent, or to update an existent Pod
func (kc KubeClient) CreateOrUpdatePod(
	ctx context.Context,
	pod *coreV1.Pod,
	mutator controllerutil.MutateFn,
) (controllerutil.OperationResult, error) {
	kc.logger.Infof("creating or updating Pod: %v/%v", pod.Namespace, pod.Name)
	return ctrl.CreateOrUpdate(ctx, kc.client, pod, mutator)
}
