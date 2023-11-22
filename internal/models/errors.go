package models

import "errors"

var (
	ErrNoRecord = errors.New("models: no matching record found")

	//Invalid creds error, used for user incorrect login
	ErrInvalidCredentials = errors.New("models: invalid credentials")

	//Duplicate email, raised when user signups with an email alrdy in use
	ErrDuplicateEmail = errors.New("models: duplicate email")
)
