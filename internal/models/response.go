package models

type Response struct {
	Code     int         `json:"code"`
	Success  bool        `json:"success"`
	Message  string      `json:"message"`
	Response interface{} `json:"response"`
}
