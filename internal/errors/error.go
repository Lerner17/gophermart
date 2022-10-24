package errors

import "errors"

var UserNameAlreadyExists = errors.New("username already exists")
var InvalidPasswordPattern = errors.New("password pattern is invalid")
var InvalidLoginOrPassword = errors.New("invalid login or password")
