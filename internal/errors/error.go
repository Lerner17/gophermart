package errors

import "errors"

var ErrUserNameAlreadyExists = errors.New("username already exists")
var ErrInvalidPasswordPattern = errors.New("password pattern is invalid")
var ErrInvalidLoginOrPassword = errors.New("invalid login or password")
var ErrOrderNumberAlreadyExists = errors.New("order number already exists")
var ErrOrderWasCreatedBySelf = errors.New("order was already created by user")
var ErrOrderWasCreatedByAnotherUser = errors.New("order was already created by another user")

var ErrOrdersNotFound = errors.New("has not orders by current user")
var ErrCannotFindOrderByNumber = errors.New("cannot find order by number")

var ErrBalanceTooLow = errors.New("user balance is too low")
