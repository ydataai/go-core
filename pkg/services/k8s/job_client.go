package k8s

import (
	"context"

	batchV1 "k8s.io/api/batch/v1"
	batchV1beta1 "k8s.io/api/batch/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// >>>>>>>>>>  JOB  >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

// GetJob gets a specific job on cluster
func (kc KubeClient) GetJob(ctx context.Context, uid string, namespace string) (*batchV1.Job, error) {
	kc.logger.Infof("fetching Job: %v/%v", namespace, uid)
	job := &batchV1.Job{}
	err := kc.client.Get(ctx, client.ObjectKey{
		Name:      uid,
		Namespace: namespace,
	}, job)

	return job, err
}

// UpdateJob creates a job on cluster
func (kc KubeClient) UpdateJob(ctx context.Context, job *batchV1.Job) error {
	kc.logger.Infof("creating or updating Job: %v/%v", job.Namespace, job.Name)
	return kc.client.Update(ctx, job)
}

// DeleteJob deletes a job on cluster
func (kc KubeClient) DeleteJob(ctx context.Context, job *batchV1.Job) error {
	kc.logger.Infof("deleting Job: %v/%v", job.Namespace, job.Name)
	return kc.client.Delete(ctx, job)
}

// >>>>>>>>>>  CRON JOB  >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

// GetCronJob gets a specific cronjob on cluster
func (kc KubeClient) GetCronJob(ctx context.Context, uid string, namespace string) (batchV1beta1.CronJob, error) {
	kc.logger.Infof("fetching CronJob: %v/%v", namespace, uid)
	cronjob := batchV1beta1.CronJob{}
	err := kc.client.Get(ctx, client.ObjectKey{
		Name:      uid,
		Namespace: namespace,
	}, &cronjob)

	return cronjob, err
}

// DeleteCronJob deletes a cronjob on cluster
func (kc KubeClient) DeleteCronJob(ctx context.Context, cronjob *batchV1beta1.CronJob) error {
	kc.logger.Infof("deleting CronJob: %v/%v", cronjob.Namespace, cronjob.Name)
	return kc.client.Delete(ctx, cronjob)
}

// CreateOrUpdateCronJob creates or updates a cronjob in the cluster
func (kc KubeClient) CreateOrUpdateCronJob(
	ctx context.Context,
	cronjob *batchV1beta1.CronJob,
	mutator controllerutil.MutateFn,
) (controllerutil.OperationResult, error) {
	kc.logger.Infof("creating or updating CronJob: %v/%v", cronjob.Namespace, cronjob.Name)

	return ctrl.CreateOrUpdate(ctx, kc.client, cronjob, mutator)
}
