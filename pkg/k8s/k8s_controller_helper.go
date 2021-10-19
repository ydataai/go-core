package k8s

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

// SetManualControllerReference sets owner as a Controller OwnerReference on controlled.
// This is used for garbage collection of the controlled object and for
// reconciling the owner object on changes to controlled (with a Watch + EnqueueRequestForOwner).
// Since only one OwnerReference can be a controller, it returns an error if
// there is another OwnerReference with Controller flag set.
func SetManualControllerReference(owner, controlled metav1.Object, scheme *runtime.Scheme) func() error {
	return func() error {
		ro, ok := owner.(runtime.Object)
		if !ok {
			return fmt.Errorf("%T is not a runtime.Object, cannot call SetControllerReference", owner)
		}
		gvk, err := apiutil.GVKForObject(ro, scheme)
		if err != nil {
			return err
		}
		ref := metav1.OwnerReference{
			APIVersion:         gvk.GroupVersion().String(),
			Kind:               gvk.Kind,
			Name:               owner.GetName(),
			UID:                owner.GetUID(),
			BlockOwnerDeletion: pointer.BoolPtr(true),
			Controller:         pointer.BoolPtr(true),
		}
		owners := controlled.GetOwnerReferences()
		if idx := indexOwnerRef(owners, ref); idx == -1 {
			owners = append(owners, ref)
		} else {
			owners[idx] = ref
		}
		controlled.SetOwnerReferences(owners)
		return nil
	}
}

// indexOwnerRef returns the index of the owner reference in the slice if found, or -1.
func indexOwnerRef(ownerReferences []metav1.OwnerReference, ref metav1.OwnerReference) int {
	for index, r := range ownerReferences {
		if referSameObject(r, ref) {
			return index
		}
	}
	return -1
}

// Returns true if a and b point to the same object.
func referSameObject(a, b metav1.OwnerReference) bool {
	aGV, err := schema.ParseGroupVersion(a.APIVersion)
	if err != nil {
		return false
	}
	bGV, err := schema.ParseGroupVersion(b.APIVersion)
	if err != nil {
		return false
	}
	return aGV.Group == bGV.Group && a.Kind == b.Kind && a.Name == b.Name
}
