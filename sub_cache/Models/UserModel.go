package Models

type User struct {
	ID        string `db:"id" json:"id"`
	Login     string `db:"login" json:"login"`
	Password  string `db:"password" json:"password"`
	SecretKey string `db:"secret_key" json:"secret_key"`
}
