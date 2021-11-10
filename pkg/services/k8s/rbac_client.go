package k8s

import (
	"context"

	rbacV1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// >>>>>>>>>>  ROLE  >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

// GetRole fetches a cluster role with the provided name in the namespace
func (kc KubeClient) GetRole(ctx context.Context, namespacedName types.NamespacedName) (rbacV1.Role, error) {
	kc.logger.Infof("fetching role %v", namespacedName)

	role := rbacV1.Role{}

	err := kc.client.Get(ctx, namespacedName, &role)

	return role, err
}

// DeleteRole deletes a role in the cluster
func (kc KubeClient) DeleteRole(ctx context.Context, role rbacV1.Role) error {
	kc.logger.Infof("deleting cluster Role %s/%s", role.Namespace, role.Name)

	if err := kc.client.Delete(ctx, &role); err != nil {
		kc.logger.Errorf("failed to delete Role: %s/%s with error %v", role.Namespace, role.Name, err)
		return err
	}

	return nil
}

// CreateOrUpdateRole creates or updates a role in the cluster
func (kc KubeClient) CreateOrUpdateRole(
	ctx context.Context,
	role *rbacV1.Role,
	mutator controllerutil.MutateFn,
) (controllerutil.OperationResult, error) {
	kc.logger.Infof("creating or updating Role: %v/%v", role.Namespace, role.Name)

	return ctrl.CreateOrUpdate(ctx, kc.client, role, mutator)
}

// >>>>>>>>>>  ROLE BINDING  >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

// GetRoleBinding fetches a role binding with the provided name in the namespace
func (kc KubeClient) GetRoleBinding(
	ctx context.Context,
	namespacedName types.NamespacedName,
) (rbacV1.RoleBinding, error) {
	kc.logger.Infof("fetching RoleBinding: %v", namespacedName)

	roleBinding := rbacV1.RoleBinding{}

	err := kc.client.Get(ctx, namespacedName, &roleBinding)

	return roleBinding, err
}

// DeleteRoleBinding deletes a role binding in the cluster
func (kc KubeClient) DeleteRoleBinding(ctx context.Context, rb rbacV1.RoleBinding) error {
	kc.logger.Infof("deleting RoleBinding: %s/%s", rb.Namespace, rb.Name)

	if err := kc.client.Delete(ctx, &rb); err != nil {
		kc.logger.Errorf("failed to delete RoleBinding: %s/%s with error %v", rb.Namespace, rb.Name, err)
		return err
	}

	return nil
}

// CreateOrUpdateRoleBinding creates or updates a role binding in the cluster
func (kc KubeClient) CreateOrUpdateRoleBinding(
	ctx context.Context,
	rb *rbacV1.RoleBinding,
	mutator controllerutil.MutateFn,
) (controllerutil.OperationResult, error) {
	kc.logger.Infof("creating or updating RoleBinding: %v/%v", rb.Namespace, rb.Name)

	return ctrl.CreateOrUpdate(ctx, kc.client, rb, mutator)
}
