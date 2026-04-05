// internal/webhook/payload.go
// Payload shapes posted to webhook endpoints.
package webhook

import "time"

// EventName is the string identifier sent in every webhook payload.
type EventName string

const (
	EventOnGenerate    EventName = "on_generate"
	EventOnRunComplete EventName = "on_run_complete"
)

// GeneratePayload is posted for each completed operation (on_generate).
type GeneratePayload struct {
	Event     EventName `json:"event"`
	Timestamp time.Time `json:"timestamp"`
	Operation struct {
		ID     string `json:"id"`
		Method string `json:"method"`
		Path   string `json:"path"`
	} `json:"operation"`
	CaseCount int `json:"case_count"`
}

// RunCompletePayload is posted once after rendering finishes (on_run_complete).
type RunCompletePayload struct {
	Event      EventName `json:"event"`
	Timestamp  time.Time `json:"timestamp"`
	TotalCases int       `json:"total_cases"`
	OutputDir  string    `json:"output_dir"`
}
