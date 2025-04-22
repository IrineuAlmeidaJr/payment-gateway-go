package domain

import "errors"

var (
	ErrAccountNotFound    = errors.New("account not found")
	ErrDuplicatedAPIKey   = errors.New("api key already exits")
	ErrInvoiceNotFound    = errors.New("invoice not found")
	ErrUnauthorizedAccess = errors.New("unauthorizes access")

	ErrInvalidAccountID = errors.New("invalid account ID")
	ErrInvalidAmount    = errors.New("invalid amount")
	ErrInvalidStatus    = errors.New("invalid status")
)
