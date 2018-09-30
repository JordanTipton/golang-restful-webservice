package errors

// NotFound error
const NotFound = Error("Not Found")

// BadRequest error
const BadRequest = Error("Bad Request")

//Error type
type Error string

func (e Error) Error() string { return string(e) }
