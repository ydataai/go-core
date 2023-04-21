// Package kubernetes is an util library to deal with kubernetes.
package kubernetes

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func createOwnerReference(
	owner,
	obj metav1.Object,
	scheme *runtime.Scheme,
) controllerutil.MutateFn {
	return func() error { return controllerutil.SetControllerReference(owner, obj, scheme) }
}
