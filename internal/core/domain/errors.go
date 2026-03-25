package domain

import "errors"

// Common Domain Errors
var (
	ErrNoRecordsReturned   = errors.New("no records returned")
	ErrIdempotencyConflict = errors.New("idempotency conflict detected")
	ErrInvalidToken        = errors.New("invalid token")
	ErrExpiredToken        = errors.New("expired token")
)

// Auth Errors
var (
	ErrRegisterUserFailed = errors.New("register user failed")
	ErrEmailAlreadyExists = errors.New("email already exists")
)
