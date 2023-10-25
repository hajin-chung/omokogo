package models

type RegisterReq struct {
	Name     string
	Password string
	Email    string
}

type LoginReq struct {
	Name     string
	Password string
}

type User struct {
	Id    string
	Name  string
	Email string
}
