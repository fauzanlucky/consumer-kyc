package constant

import (
	"errors"
)

var (
	ErrConsumerPrefix = errors.New("consumer error")

	ErrNoAuthorizationFound = errors.New("no authorization found")
	ErrKYCAlreadyHandled    = errors.New("KYC request already handled")
	ErrInvalidID            = errors.New("invalid id")
	ErrSimilarIDNotFound    = errors.New("similar account id not found")
)
