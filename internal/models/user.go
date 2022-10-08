package models

type RegisterUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
