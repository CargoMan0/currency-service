package errors

import "errors"

// External errors
var (
	ErrClientUnexpectedStatusCode = errors.New("client: unexpected status code")
	ErrClientInvalidCredentials   = errors.New("client: invalid credentials")
	ErrClientTokenGeneration      = errors.New("client: token generation failed")

	ErrClientTokenNotFound         = errors.New("client: token not found in header")
	ErrClientInvalidOrExpiredToken = errors.New("client: invalid signature or token expired")
)
