package gateway_b

import (
	"bytes"
	"context"
	"encoding/xml"
)

func (c ClientWithResponses) Call(
	ctx context.Context,
	txnID string,
	amount float32,
	currency string,
) error {
	xmlData, err := xml.MarshalIndent(
		struct {
			XMLName       xml.Name `xml:"Transaction"`
			TransactionID string   `xml:"TransactionID"`
			Amount        float32  `xml:"Amount"`
			Currency      string   `xml:"Currency"`
		}{
			TransactionID: txnID,
			Amount:        amount,
			Currency:      currency,
		}, "", "  ")
	if err != nil {
		return err
	}

	if _, err := c.PostProcessPaymentWithBodyWithResponse(ctx, "application/xml", bytes.NewReader(xmlData), nil); err != nil {
		return err
	}

	return nil
}
