# Payment Processing Service (WIP)

**Note:** This solution is **not production-ready** and is still under development. Several critical features are currently missing, including:

- The callback mechanism is not yet implemented.
- The withdrawal functionality is incomplete.
- Other requirements are also pending implementation.

## Focus of This Solution

This solution emphasizes resilience and reliability by using the **Outbox Pattern** combined with a **worker pool**. This approach ensures that no information is lost and that interactions with third-party services are handled asynchronously, separate from request processing. As a result, the application is more resilient to sudden unavailability, which is common in cloud environments.

## Key Areas to Focus On

### 1. Transactional Operations in Deposit Use Case
File: `exinity/internal/usecases/deposit/usecase.go`

In this file, the deposit operation is implemented transactionally. The process includes recording the transaction details in the database and creating a task in the outbox, ensuring both actions are atomic and consistent.

### 2. Client Code for Third-Party Gateway
File: `exinity/internal/clients/gateway_a/client.go`

The client code for interacting with a third-party payment gateway is generated and includes retry mechanisms to handle transient errors, ensuring that the system is resilient to temporary failures.

### 3. Outbox Worker Implementation
File: `exinity/internal/outbox/jobs/deposit/job.go`

The worker handling the outbox tasks also includes retry logic and can implement **circuit breaker** patterns. This makes sure that the processing of deposit tasks is robust, with mechanisms to prevent system overload in case of repeated failures.

## Testing and Code Generation

- **Tests**: Some tests are implemented in the `exinity/internal/clients/gateway_a/` directory.
    - To run tests, use:
      ```bash
      go test ./...
      ```

- **Code Generation**: The client code for third-party interactions is generated.
    - To regenerate the code, use:
      ```bash
      go generate ./...
      ```
