package errors

// NotFound error
const NotFound = Error("sql: no rows in result set")

//Error type
type Error string

func (e Error) Error() string { return string(e) }
