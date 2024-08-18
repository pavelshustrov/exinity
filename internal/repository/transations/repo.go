package transations

import (
	"context"
	"exinity/internal/database"
)

type Repo struct {
	dbClient *database.Transactor
}

func New(dbClient *database.Transactor) *Repo {
	return &Repo{dbClient: dbClient}
}

func (repo *Repo) Create(
	ctx context.Context,
	txnType string,
	amount float32,
	currencyCode string,
	paymentGateway string,
) (string, error) {
	query := `
    INSERT INTO txn_details (txn_type, amount, currency_code, payment_gateway)
    VALUES ($1, $2, $3, $4) RETURNING txn_id;
    `

	var txnID string
	err := repo.dbClient.LoadClient(ctx).QueryRow(query, txnType, amount, currencyCode, paymentGateway).Scan(&txnID)
	if err != nil {
		return "", err
	}

	return txnID, nil
}

func (repo *Repo) UpdateStatus(
	ctx context.Context,
	txnID string,
	status string,
	details *string,
) error {
	query := `
    INSERT INTO txn_statuses (txn_id, new_status, details)
    VALUES ($1, $2, $3, $4);
    `

	_, err := repo.dbClient.LoadClient(ctx).Exec(query, txnID, status, details)
	if err != nil {
		return err
	}

	return nil
}
