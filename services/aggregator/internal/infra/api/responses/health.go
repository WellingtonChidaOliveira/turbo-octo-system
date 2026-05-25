package responses

type Health struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}
