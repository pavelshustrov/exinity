package deposit

import "exinity/internal/validator"

type Request struct {
	OrderID  string  `validate:"required,uuid"`
	Amount   float32 `validate:"required"`
	Currency string  `validate:"required,currency"`
	Gateway  string  `validate:"required"`
}

func (r Request) Validate() error {
	return validator.Validator.Struct(r)
}
