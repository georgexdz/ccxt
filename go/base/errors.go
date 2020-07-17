package base

import (
	"errors"
	"fmt"
)

var (
	BaseError = errors.New("")

	InternalError = fmt.Errorf("%w", BaseError)

	ExchangeError = fmt.Errorf("%w", BaseError)

	AuthenticationError = fmt.Errorf("%w", ExchangeError)
	PermissionDenied = fmt.Errorf("%w", AuthenticationError)
	AccountSuspended = fmt.Errorf("%w", AuthenticationError)

	ArgumentsRequired = fmt.Errorf("%w", ExchangeError)
	BadRequest = fmt.Errorf("%w", ExchangeError)
	BadSymbol = fmt.Errorf("%w", BadRequest)

	BadResponse = fmt.Errorf("%w", ExchangeError)
	NullResponse = fmt.Errorf("%w", BadResponse)

	InsufficientFunds = fmt.Errorf("%w", ExchangeError)
	InvalidAddress = fmt.Errorf("%w", ExchangeError)
	AddressPending = fmt.Errorf("%w", InvalidAddress)

	InvalidOrder = fmt.Errorf("%w", ExchangeError)
	OrderNotFound = fmt.Errorf("%w", InvalidOrder)
	OrderNotCached = fmt.Errorf("%w", InvalidOrder)
	CancelPending = fmt.Errorf("%w", InvalidOrder)
	OrderImmediatelyFillable = fmt.Errorf("%w", InvalidOrder)
	OrderNotFillable = fmt.Errorf("%w", InvalidOrder)
	DuplicateOrderId = fmt.Errorf("%w", InvalidOrder)

	NotSupported = fmt.Errorf("%w", ExchangeError)

	NetworkError = fmt.Errorf("%w", BaseError)
	DDoSProtection = fmt.Errorf("%w", NetworkError)
	RateLimitExceeded = fmt.Errorf("%w", DDoSProtection)
	ExchangeNotAvailable = fmt.Errorf("%w", NetworkError)
	OnMaintenance = fmt.Errorf("%w", ExchangeNotAvailable)
	InvalidNonce = fmt.Errorf("%w", NetworkError)
	RequestTimeout = fmt.Errorf("%w", NetworkError)
)


// TypedError creates a typed error from type t and message, if type does
// not match a known error type, fmt.Errorf will be used
func TypedError(t string, msg string) error {
	var err error
	switch t {
	case "BaseError":
		err = BaseError
	case "AccountSuspended":
		err = AccountSuspended
	case "NullResponse":
		err = NullResponse
	case "AddressPending":
		err = AddressPending
	case "OrderImmediatelyFillable":
		err = OrderImmediatelyFillable
	case "OrderNotFillable":
		err = OrderNotFillable
	case "DuplicateOrderId":
		err = DuplicateOrderId
	case "OnMaintenance":
		err = OnMaintenance
	case "NotSupported":
		err = NotSupported
	case "RateLimitExceeded":
		err = RateLimitExceeded
	case "ArgumentsRequired":
		err = ArgumentsRequired
	case "AuthenticationError":
		err = AuthenticationError
	case "InvalidNonce":
		err = InvalidNonce
	case "InsufficientFunds":
		err = InsufficientFunds
	case "InvalidOrder":
		err = InvalidOrder
	case "OrderNotFound":
		err = OrderNotFound
	case "OrderNotCached":
		err = OrderNotCached
	case "PermissionDenied":
		err = PermissionDenied
	case "CancelPending":
		err = CancelPending
	case "NetworkError":
		err = NetworkError
	case "DDoSProtection":
		err = DDoSProtection
	case "RequestTimeout":
		err = RequestTimeout
	case "ExchangeNotAvailable":
		err = ExchangeNotAvailable
	case "BadSymbol":
		err = BadSymbol
	case "InternalError":
		err = InternalError
	default:
		err = errors.New(t)
	}

	return fmt.Errorf("%w%v:%v", err, t, msg)
}

