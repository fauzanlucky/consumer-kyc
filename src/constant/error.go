package constant

import (
	"errors"
)

var (
	ErrNoAuthorizationFound = errors.New("no authorization found")
	ErrKYCAlreadyHandled    = errors.New("KYC request already handled")
	ErrInvalidID            = errors.New("invalid id")
)
