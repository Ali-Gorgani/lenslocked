package errors

import (
	"errors"
)

// Provides access to the standard error handling functions "Is" and "As" from the "errors" package.
// "Is" checks if an error matches a specific error in a chain of wrapped errors.
// "As" attempts to cast an error to a specified type, allowing for type assertions.
var (
	Is = errors.Is
	As = errors.As
)