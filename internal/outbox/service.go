package outbox

import (
	"context"
	"errors"
	"exinity/internal/repository/jobs"
	"github.com/labstack/gommon/log"
	"time"
)

type jobRepository interface {
	Create(ctx context.Context, name string, payload string) (int, error)
	GetPending(ctx context.Context) (jobs.Job, error)
	Complete(ctx context.Context, jobID int) error
	Fail(ctx context.Context, jobID int) (int, error)
}

type Service struct {
	jobRepo     jobRepository
	idleTime    time.Duration
	jobRegistry map[string]Job
}

func NewService(jobRepo jobRepository, idleTime time.Duration) *Service {
	return &Service{jobRepo: jobRepo, idleTime: idleTime, jobRegistry: make(map[string]Job)}
}

func (s *Service) RegisterJob(name string, job Job) {
	s.jobRegistry[name] = job
}

func (s *Service) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		job, err := s.jobRepo.GetPending(ctx)
		if err != nil {
			switch {
			case errors.Is(err, context.Canceled):
				return nil
			case errors.Is(err, jobs.ErrNoJobs):
				select {
				case <-ctx.Done():
					return nil
				case <-time.After(s.idleTime):
					continue
				}
			default:
				log.Debug(
					"error on find and reserve job",
					"job.ID",
					job.ID,
				)

				return err
			}
		}

		if err := s.work(ctx, job); err != nil {
			if _, err := s.jobRepo.Fail(ctx, job.ID); err != nil {
				log.Errorf("error on fail job: %v", err)
			}
			return err
		}

		if err := s.jobRepo.Complete(ctx, job.ID); err != nil {
			return err
		}
	}
}

func (s *Service) work(ctx context.Context, job jobs.Job) error {
	return s.jobRegistry[job.EventType].Handle(ctx, job.Payload)
}
