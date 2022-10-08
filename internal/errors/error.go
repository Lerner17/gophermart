package errors

import "errors"

var UserNameAlreadyExists = errors.New("username already exists")
var InvalidPasswordPattern = errors.New("password pattern is invalid")
