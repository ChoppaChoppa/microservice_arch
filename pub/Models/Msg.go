package Models

type Message struct {
	Order     OrderInfo `json:"order"`
	SecretKey string    `json:"secret_key"`
}
