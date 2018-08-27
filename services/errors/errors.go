package errors

// NotFound error
const NotFound = Error("Not Found")

//Error type
type Error string

func (e Error) Error() string { return string(e) }
