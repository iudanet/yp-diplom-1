package models

type User struct {
	Balanse int64  `json:"balanse"`
	Login   string `json:"login"`
}

type UserAuth struct {
	ID           int64  `json:"id"`
	Login        string `json:"login"`
	PasswordHash string `json:"-"`
}
