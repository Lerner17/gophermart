package errors

import "fmt"

type HTTPError struct {
	Code int
	Msg  string
}

func (e HTTPError) HTTPCode() int {
	return e.Code
}
func (e HTTPError) Error() string {
	return fmt.Sprintf("HTTPError %v: %s", e.Code, e.Msg)
}
func (e HTTPError) Message() string {
	return e.Msg
}
