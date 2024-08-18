package gateway_a

import (
	"context"
	"errors"
)

func (c ClientWithResponses) Call(
	ctx context.Context,
	txnID string,
	amount float32,
	currency string,
) error {
	resp, err := c.PostPaymentsWithResponse(ctx, PostPaymentsJSONRequestBody{
		Amount:        amount,
		Currency:      currency,
		Type:          "deposit",
		TransactionId: txnID,
	})

	if err != nil {
		return err
	}

	// Check if the status code is not 2xx and return an error
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return errors.New("received non-2xx response from API")
	}

	return nil
}
