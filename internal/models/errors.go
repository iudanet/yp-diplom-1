package models

import "errors"

var (
	ErrInvalidOrderNumber                = errors.New("invalid order number")
	ErrOrderAlreadyUploaded              = errors.New("order already uploaded by this user")
	ErrOrderAlreadyUploadedByAnotherUser = errors.New("order already uploaded by another user")
	ErrUserNotFound                      = errors.New("user not found")
	ErrInsufficientFunds                 = errors.New("insufficient funds")
)
