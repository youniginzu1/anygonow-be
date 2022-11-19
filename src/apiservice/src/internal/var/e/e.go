package e

import (
	"golang.org/x/xerrors"
)

var (
	ErrMissingField         = func(a string) error { return xerrors.Errorf("missing field: %s", a) }
	ErrBodyInvalid          = xerrors.New("body invalid")
	ErrMissingBody          = xerrors.New("missing body")
	ErrMissingCertificate   = xerrors.New("missing certificate")
	ErrExpectedBody         = xerrors.New("expected body")
	ErrOTPFail              = xerrors.New("otp fail")
	ErrIdInvalidFormat      = xerrors.New("invalid id")
	ErrNoPermission         = xerrors.New("no permisssion")
	ErrInvalidFile          = xerrors.New("invalid file")
	ErrExceedMaxOrders      = xerrors.New("exceed max orders")
	ErrAlreadyOrdered       = xerrors.New("already ordered")
	ErrInvalidOrderStatus   = xerrors.New("invalid order status")
	ErrCategoryExisted      = xerrors.New("service existed")
	ErrPayment              = xerrors.New("payment error")
	ErrGroupExisted         = xerrors.New("group existed")
	ErrFeedbackExisted      = xerrors.New("you already reviewed this order")
	RefCodeNotFound         = xerrors.New("cannot found your referral account")
	ExceedBuyAdvertiseLimit = xerrors.New("exceed limit for buying advertise")
	ErrStripeHeader         = xerrors.New("cannot get stripe header")
)
