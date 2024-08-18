package deposit

import (
	"context"
	"encoding/json"
)

var Name = "deposit"

type Job struct {
	callersMap map[string]Caller
}

func NewJob(callers map[string]Caller) *Job {
	return &Job{callersMap: make(map[string]Caller)}

}
func (j Job) Name() string {
	return Name
}

type Caller interface {
	Call(
		ctx context.Context,
		txnID string,
		amount float32,
		currency string,
	) error
}

func (j Job) Handle(ctx context.Context, data string) error {
	var payload Payload

	if err := json.Unmarshal([]byte(data), &payload); err != nil {
		return err
	}

	caller := j.callersMap[payload.Gateway]
	return caller.Call(ctx, payload.TxnID, payload.Amount, payload.Currency)
}
