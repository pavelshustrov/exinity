package main

import (
	"context"
	"database/sql"
	"exinity/internal/clients"
	"exinity/internal/clients/gateway_a"
	"exinity/internal/clients/gateway_b"
	"exinity/internal/conterollers"
	"exinity/internal/database"
	"exinity/internal/outbox"
	"exinity/internal/outbox/jobs/deposit"
	"exinity/internal/repository/jobs"
	"exinity/internal/repository/transations"
	"exinity/internal/server"
	depositusecase "exinity/internal/usecases/deposit"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Hello World")
	httpClient := &http.Client{
		Timeout:   10 * time.Second, // Set the timeout for the request
		Transport: &clients.CustomTransport{Transport: http.DefaultTransport},
	}

	gatewayAClient, err := gateway_a.NewClientWithResponses("gateway_a.com", gateway_a.WithHTTPClient(httpClient))
	if err != nil {
		log.Fatal("Failed to create client for gateway_a:", err)
	}

	gatewayBClient, err := gateway_b.NewClientWithResponses("gateway_b.com", gateway_b.WithHTTPClient(httpClient))
	if err != nil {
		log.Fatal("Failed to create client for gateway_b:", err)
	}

	depositJob := deposit.NewJob(map[string]deposit.Caller{
		"gateway_a": gatewayAClient,
		"gateway_b": gatewayBClient,
	})

	connStr := "user=username dbname=yourdb sslmode=disable password=yourpassword host=localhost port=5432"

	// Open a connection to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open a DB connection:", err)
	}
	defer db.Close()

	transactor := database.NewTransactor(db)
	jobsRepo := jobs.New(transactor)
	txnRepo := transations.New(transactor)

	outboxSrv := outbox.NewService(jobsRepo, 500*time.Millisecond)
	outboxSrv.RegisterJob("deposit", depositJob)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	go func() {
		if err := outboxSrv.Run(ctx); err != nil {
			log.Println("Failed to start outbox:", err)
		}
	}()

	depositUsecase := depositusecase.New(transactor, txnRepo, jobsRepo)

	handler := conterollers.New(depositUsecase)

	e := echo.New()
	e.Use(
		middleware.Logger(),
		middleware.Recover(),
	)

	server.RegisterHandlers(e, handler)

	srv := &http.Server{
		Addr:    "localhost:8080",
		Handler: e,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
