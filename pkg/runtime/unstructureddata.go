package runtime

import "k8s.io/apimachinery/pkg/runtime"

type UnstructuredData map[string]interface{}

func (in *UnstructuredData) DeepCopy() *UnstructuredData {
	if in == nil {
		return nil
	}

	out := *in
	out = runtime.DeepCopyJSON(out)
	return &out
}
