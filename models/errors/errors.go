package errors

// NotFound error type
type NotFound struct {
	Message string
}

// InvalidArgument error type
type InvalidArgument struct {
	Message string
}

// Error method for NotFound
func (e NotFound) Error() string { return e.Message }

// Error method for InvalidArgument
func (e InvalidArgument) Error() string { return e.Message }
