package maas_config

import "errors"

var (
	ErrInvalidAuthorityForCaller = errors.New("invalid authority for caller")
	ErrInvalidAuthorityForOwner = errors.New("invalid authority for owner")
	ErrInvalidInput = errors.New("invalid input")
	ErrEmitEventChangeOwner = errors.New("emit EventChangeOwner error")
	ErrEncodeValueFailed = errors.New("encode value failed")
	ErrEmitBlockAccount = errors.New("emitBlockAccount error")
)

