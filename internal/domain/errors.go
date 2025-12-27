package domain

import "errors"

var (
	ErrNotFound               = errors.New("resource not found")
	ErrEmailAlreadyRegistered = errors.New("данный email уже зарегистрирован")
	ErrInvalidCredentials     = errors.New("неверный логин или пароль")
)
