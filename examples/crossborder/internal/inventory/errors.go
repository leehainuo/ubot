package inventory

import "errors"

var (
	ErrNotFound = errors.New(
		"inventory not found",
	)

	ErrInvalidTransfer = errors.New(
		"invalid inventory transfer",
	)

	ErrInsufficientStock = errors.New(
		"insufficient inventory",
	)

	ErrIdempotencyKey = errors.New(
		"idempotency key is required",
	)

	ErrIdempotencyConflict = errors.New(
		"idempotency key was already used for different request",
	)
)
