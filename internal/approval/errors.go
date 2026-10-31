package approval

import "errors"

var (
	ErrRequestNotFound = errors.New("request not found")
	ErrAlreadyResolved = errors.New("request already resolved")
	ErrCancelled       = errors.New("request cancelled")
	ErrExpired         = errors.New("request expired")
	ErrNotAllowed      = errors.New("not allowed to resolve this request")
	ErrInvalidDecision = errors.New("decision must be approved or rejected")
)
