package message

type Message struct {
	Type string `json:"type"`
	body []byte `json:"body"`
}
