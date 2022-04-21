package entity

type Message struct {
	Default string `json:"default"`
}

type EmailTokenMessage struct {
	Email       string `json:"email"`
	Token       string `json:"token"`
	MessageType string `json:"message_type"`
}
