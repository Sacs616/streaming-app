package domain

import "errors"

var (
	ErrVideoNotFound      = errors.New("video not found")
	ErrInvalidVideoID     = errors.New("invalid video ID")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrVideoAlreadyExists = errors.New("video already exists")
)
