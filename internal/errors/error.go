package errors

import "errors"

var UserNameAlreadyExists = errors.New("username already exists")
var InvalidPasswordPattern = errors.New("password pattern is invalid")
var InvalidLoginOrPassword = errors.New("invalid login or password")
var OrderNumberAlreadyExists = errors.New("order number already exists")
var OrderWasCreatedBySelf = errors.New("order was already created by user")
var OrderWasCreatedByAnotherUser = errors.New("order was already created by another user")

var OrdersNotFound = errors.New("has not orders by current user")
var CannotFindOrderByNumber = errors.New("cannot find order by number")
