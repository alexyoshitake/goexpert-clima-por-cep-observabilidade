package domain

import "errors"

var (
	ErrCEPInvalido      = errors.New("invalid zipcode")
	ErrCEPNaoEncontrado = errors.New("can not find zipcode")
)
