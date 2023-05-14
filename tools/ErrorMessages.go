package tools

import "errors"

// ERROR_WRONG_NUMBER_ARGUMENTS returns the error message if the wrong number of arguments have been provided
const ERROR_WRONG_NUMBER_ARGUMENTS = "Error: You did not provide a valid amount of arguments."

// CreateAlphaNumericError returns an error object with the hint that the name of the object passed by a string should be
// alphanumeric.
func CreateAlphaNumericError(objectName string) error {
	return errors.New(objectName + ": Name has to be alphanumeric.")
}
