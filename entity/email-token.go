package entity

type EmailToken struct {
	Email string
	Token string
	TTL   int64
}
