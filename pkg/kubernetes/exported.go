// Package kubernetes is an util library to deal with kubernetes.
package kubernetes

// SetOwnerReference sets owner as a Controller OwnerReference on controlled.
// This is used for garbage collection of the controlled object and for
// reconciling the owner object on changes to controlled (with a Watch + EnqueueRequestForOwner).
// Since only one OwnerReference can be a controller, it returns an error if
// there is another OwnerReference with Controller flag set.
var SetOwnerReference = setOwnerReference

// SetCrossNamespaceOwnerReference allows you to set an owner from another namespace
// This is used for garbage collection of the controlled object and for
// reconciling the owner object on changes to controlled (with a Watch + EnqueueRequestForOwner).
// Since only one OwnerReference can be a controller, it returns an error if
// there is another OwnerReference with Controller flag set
var SetCrossNamespaceOwnerReference = setCrossNamespaceOwnerReference
