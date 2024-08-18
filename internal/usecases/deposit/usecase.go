package deposit

import (
	"context"
	"encoding/json"
	"errors"
	"exinity/internal/outbox/jobs/deposit"
	"fmt"
)

var (
	ErrInvalidRequest = errors.New("invalid request")
)

type transactor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) (err error)
}

type txnRepo interface {
	Create(
		ctx context.Context,
		txnType string,
		amount float32,
		currencyCode string,
		paymentGateway string,
	) (string, error)
	UpdateStatus(
		ctx context.Context,
		txnID string,
		status string,
		details *string,
	) error
}

type jobRepository interface {
	Create(
		ctx context.Context,
		name string,
		payload string,
	) (int, error)
}

type Usecase struct {
	txn  transactor
	repo txnRepo

	jobRepo jobRepository
}

func New(trct transactor, txnRepo txnRepo, jobRepo jobRepository) *Usecase {
	return &Usecase{
		txn:     trct,
		repo:    txnRepo,
		jobRepo: jobRepo,
	}
}

func (uc *Usecase) Handle(ctx context.Context, req Request) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("%w:%v", ErrInvalidRequest, err)
	}

	return uc.txn.RunInTx(ctx, func(ctx context.Context) error {
		txnID, err := uc.repo.Create(ctx, "deposit", req.Amount, req.Currency, req.Gateway)
		if err != nil {
			return err
		}
		if err := uc.repo.UpdateStatus(ctx, txnID, "pending", nil); err != nil {
			return err
		}

		bytes, err := json.Marshal(deposit.Payload{
			Amount:   req.Amount,
			Currency: req.Currency,
			Gateway:  req.Gateway,
			TxnID:    txnID,
		})
		if err != nil {
			return err
		}

		if _, err := uc.jobRepo.Create(ctx, deposit.Name, string(bytes)); err != nil {
			return err
		}

		return nil
	})
}
