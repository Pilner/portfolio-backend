package domain

import "errors"

// Common Domain Errors
var (
	ErrNoRecordsReturned         = errors.New("no records returned")
	ErrIdempotencyConflict       = errors.New("idempotency conflict detected")
	ErrInvalidToken              = errors.New("invalid token")
	ErrExpiredToken              = errors.New("expired token")
	ErrInconsistentDataRelations = errors.New("inconsistent data relations")
)

// Auth Errors
var (
	ErrCreateUserFailed     = errors.New("create user failed")
	ErrEmailAlreadyExists   = errors.New("email already exists")
	ErrFindUserFailed       = errors.New("find user failed")
	ErrPasswordDoesNotMatch = errors.New("password does not match")
)
