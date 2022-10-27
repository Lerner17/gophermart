package models

type User struct {
	ID       int    `json:"user_id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
