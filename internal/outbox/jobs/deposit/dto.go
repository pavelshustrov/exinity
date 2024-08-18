package deposit

type Payload struct {
	Amount   float32 `json:"amount" validate:"required"`
	Currency string  `json:"currency" validate:"required,currency"`
	Gateway  string  `json:"gateway" validate:"required"`
	TxnID    string  `json:"txn_id" validate:"required,uuid"`
}
