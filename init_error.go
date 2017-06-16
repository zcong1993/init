package main

import "fmt"

// FAILED_TO_DOWNLOAD_FILE is error code of download error
const FAILED_TO_DOWNLOAD_FILE = 500

// InitError a custom error type we defined so that we can provide friendlier error messages
type InitError struct {
	errorCode int    // an error code is an arbitrary int that allows for strongly typed identification of specific errors
	details   string // the output of the underlying error message, if any
	err       error  // the underlying golang error, if any
}

// Error is the golang error interface we implemented
func (e *InitError) Error() string {
	return fmt.Sprintf("%d - %s", e.errorCode, e.details)
}

func newError(errorCode int, details string) *InitError {
	return &InitError{
		errorCode: errorCode,
		details:   details,
		err:       nil,
	}
}

func wrapError(err error) *InitError {
	if err == nil {
		return nil
	}
	return &InitError{
		errorCode: -1,
		details:   err.Error(),
		err:       err,
	}
}
