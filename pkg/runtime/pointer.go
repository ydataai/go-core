// Package runtime is a package with utils to be used in runtime execution
package runtime

func Pointer[T any](x T) *T {
	return &x
}
