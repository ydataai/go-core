package config

// UID defines UID type represents a ID for a specific resource
type UID string

// Namespace defines the namespace type that represents the cluster namespace
type Namespace string

// UserID defines the user id type
type UserID string

// Credentials store the credentials from vault
type Credentials map[string]string

// ResourceKind defines a resource kind type
type ResourceKind string
