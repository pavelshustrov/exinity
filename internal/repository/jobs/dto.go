package jobs

type Job struct {
	ID         int
	EventType  string
	Payload    string
	RetryCount int
}
