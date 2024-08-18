package conterollers

import (
	"context"
	"exinity/internal/usecases/deposit"
)

type Depositer interface {
	Handle(ctx context.Context, req deposit.Request) error
}

type Handler struct {
	depositer Depositer
}

func New(depositer Depositer) Handler {
	return Handler{
		depositer: depositer,
	}
}
