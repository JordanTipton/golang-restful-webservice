package errors

// SQLNotFound error string returned by sql driver
const SQLNotFound = "sql: no rows in result set"

// NotFound error type
type NotFound struct {
	Message string
}

// Error method for NotFound
func (e NotFound) Error() string { return e.Message }
