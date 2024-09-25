package model

type LoginResponse struct {
	UserId   string `json:"userId"`
	Roles    string `json:"roles"`
	UserName string `json:"userName"`
	Password string `json:"hashedPassword"`
	Expired  int64  `json:"expired"`
}