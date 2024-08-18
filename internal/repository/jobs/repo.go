package jobs

import (
	"context"
	"database/sql"
	"errors"
	"exinity/internal/database"
	"time"
)

var ErrNoJobs = errors.New("no jobs found")

type Repo struct {
	dbClient *database.Transactor
}

func New(dbClient *database.Transactor) *Repo {
	return &Repo{dbClient: dbClient}
}

func (repo *Repo) Create(
	ctx context.Context,
	name string,
	payload string,
) (int, error) {
	query := `
    INSERT INTO outbox (event_type, payload, status) 
              VALUES ($1, $2, 'pending') 
              RETURNING id;
    `

	var jobID int
	err := repo.dbClient.LoadClient(ctx).QueryRow(query, name, payload, time.Now()).Scan(&jobID)
	if err != nil {
		return 0, err
	}

	return jobID, nil
}

func (repo *Repo) GetPending(
	ctx context.Context,
) (Job, error) {
	query := `
    SELECT id, event_type, payload, retry_count 
              FROM outbox 
              WHERE status = 'pending' AND retry_count < 3 AND deleted_at IS NULL
              ORDER BY created_at 
              LIMIT 1;
    `

	var job Job
	err := repo.dbClient.LoadClient(ctx).
		QueryRow(query).
		Scan(&job.ID, &job.EventType, &job.Payload, &job.RetryCount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Job{}, ErrNoJobs
		}

		return Job{}, err
	}

	return job, nil
}

func (repo *Repo) Complete(
	ctx context.Context,
	jobID int,
) error {
	query := `UPDATE outbox 
              SET deleted_at = NOW() 
              WHERE id = $1`
	_, err := repo.dbClient.LoadClient(ctx).Exec(query, jobID)
	return err
}

func (repo *Repo) Fail(ctx context.Context, jobID int) (int, error) {
	var retryCount int
	query := `UPDATE outbox 
              SET retry_count = retry_count + 1, updated_at = NOW() 
              WHERE id = $1 AND deleted_at IS NULL
              RETURNING retry_count`
	err := repo.dbClient.LoadClient(ctx).QueryRow(query, jobID).Scan(&retryCount)
	return retryCount, err
}
