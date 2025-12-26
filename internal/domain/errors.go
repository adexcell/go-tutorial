package domain

import "errors"

var ErrNotFound = errors.New("resource not found")
var ErrEmailAlreadyRegistered = errors.New("данный email уже зарегистрирован")
