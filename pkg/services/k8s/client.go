package k8s

import (
	"context"

	"github.com/sirupsen/logrus"

	batchV1 "k8s.io/api/batch/v1"
	batchV1beta1 "k8s.io/api/batch/v1beta1"
	coreV1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// KubeClient implements structure for Kubernetes Client
type KubeClient struct {
	client client.Client
	logger *logrus.Logger
}

// NewKubeClient initializes a new client for Kubernetes
func NewKubeClient(logger *logrus.Logger, client client.Client) KubeClient {
	return KubeClient{
		client: client,
		logger: logger,
	}
}

// KubeClientInterface defines the interface for Kubernetes Client
type KubeClientInterface interface {
	GetPod(ctx context.Context, namespacedName types.NamespacedName) (coreV1.Pod, error)
	DeletePod(ctx context.Context, pod *coreV1.Pod) error
	CreateOrUpdatePod(
		ctx context.Context,
		pod *coreV1.Pod,
		mutator controllerutil.MutateFn,
	) (controllerutil.OperationResult, error)

	GetService(ctx context.Context, namespacedName types.NamespacedName) (coreV1.Service, error)
	DeleteService(ctx context.Context, pod *coreV1.Service) error
	CreateOrUpdateService(
		ctx context.Context,
		service *coreV1.Service,
		mutator controllerutil.MutateFn,
	) (controllerutil.OperationResult, error)

	GetPersistentVolume(ctx context.Context, namespacedName types.NamespacedName) (coreV1.PersistentVolume, error)
	DeletePersistentVolume(ctx context.Context, pvc *coreV1.PersistentVolumeClaim) error
	CreateOrUpdatePersistentVolume(
		ctx context.Context,
		pv *coreV1.PersistentVolume,
		mutator controllerutil.MutateFn,
	) (controllerutil.OperationResult, error)

	GetPersistentVolumeClaim(
		ctx context.Context,
		namespacedName types.NamespacedName,
	) (coreV1.PersistentVolumeClaim, error)
	DeletePersistentVolumeClaim(ctx context.Context, pvc *coreV1.PersistentVolumeClaim) error
	CreateOrUpdatePersistentVolumeClaim(
		ctx context.Context,
		pvc *coreV1.PersistentVolumeClaim,
		mutator controllerutil.MutateFn,
	) (controllerutil.OperationResult, error)

	GetServiceAccount(ctx context.Context, namespacedName types.NamespacedName) (coreV1.ServiceAccount, error)
	DeleteServiceAccount(ctx context.Context, sa *coreV1.ServiceAccount) error
	CreateOrUpdateServiceAccount(
		ctx context.Context,
		sa *coreV1.ServiceAccount,
		mutator controllerutil.MutateFn,
	) (controllerutil.OperationResult, error)

	GetRole(ctx context.Context, namespacedName types.NamespacedName) (rbacV1.Role, error)
	DeleteRole(ctx context.Context, role *rbacV1.Role) error
	CreateOrUpdateRole(
		ctx context.Context,
		role *rbacV1.Role,
		mutator controllerutil.MutateFn,
	) (controllerutil.OperationResult, error)

	GetRoleBinding(ctx context.Context, namespacedName types.NamespacedName) (rbacV1.RoleBinding, error)
	DeleteRoleBinding(ctx context.Context, rb *rbacV1.RoleBinding) error
	CreateOrUpdateRoleBinding(
		ctx context.Context,
		rb *rbacV1.RoleBinding,
		mutator controllerutil.MutateFn,
	) (controllerutil.OperationResult, error)

	GetJob(ctx context.Context, uid string, namespace string) (batchV1beta1.CronJob, error)
	UpdateJob(ctx context.Context, job *batchV1.Job) error
	DeleteJob(ctx context.Context, job *batchV1.Job) error

	GetCronJob(ctx context.Context, uid string, namespace string) (batchV1beta1.CronJob, error)
	DeleteCronJob(ctx context.Context, cronjob *batchV1beta1.CronJob) error
	CreateOrUpdateCronJob(
		ctx context.Context,
		cronjob *batchV1beta1.CronJob,
		mutator controllerutil.MutateFn,
	) (controllerutil.OperationResult, error)
}
