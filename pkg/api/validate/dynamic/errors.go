package dynamic

import "fmt"

// Generic errors to be returned by non-RP context below

type GenericError struct {
	ResourceID   string
	ResourceType string
	Message      string
}

// PermissionError represents a permission error on a resource
type PermissionError struct {
	*GenericError
}

func (e *PermissionError) Error() string {
	return fmt.Sprintf("%s '%s' does not have the correct permissions. %s", e.ResourceType, e.ResourceID, e.Message)
}

// NotFoundError represents a resource not found
type NotFoundError struct {
	*GenericError
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s '%s' not found", e.ResourceType, e.ResourceID)
}

// InvalidResourceError represents a resource in an invalid state for ARO usage
type InvalidResourceError struct {
	*GenericError
}

func (e *InvalidResourceError) Error() string {
	return fmt.Sprintf("%s '%s' has attributes that make it invalid: %s", e.ResourceType, e.ResourceID, e.Message)
}

// InvalidCredentialsError represents invalid credentials attempted to be used
type InvalidCredentialsError struct {
	*GenericError
}

func (*InvalidCredentialsError) Error() string {
	return "The provided service principal credentials are invalid"
}

// InvalidTokenClaims represents a token returned which does not contain the required claims
type InvalidTokenClaims struct {
	*GenericError
}

func (*InvalidTokenClaims) Error() string {
	return "The token does not contain the required claims"
}
