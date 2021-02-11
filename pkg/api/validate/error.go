package validate

import "fmt"

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

// PermissionError represents a permission error on a resource
type PermissionError struct {
	resourceID   string
	resourceType string
}

func (e *PermissionError) Error() string {
	return fmt.Sprintf("%s '%s' does not have the correct permissions", e.resourceType, e.resourceID)
}

// NotFoundError represents a resource not found
type NotFoundError struct {
	resourceID   string
	resourceType string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s '%s' not found", e.resourceType, e.resourceID)
}

// InvalidResourceError represents a resource in an invalid state for ARO usage
type InvalidResourceError struct {
	resourceID   string
	resourceType string
	message      string
}

func (e *InvalidResourceError) Error() string {
	return fmt.Sprintf("%s '%s' has attributes that make it invalid: %s", e.resourceType, e.resourceID, e.message)
}
