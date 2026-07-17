package bark

type BarkPayload struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body"`
}
