package domain

type User struct {
	Id          string
	Email       string
	DisplayName string
}

type AuthBase struct {
	Email    string
	Password string
}

type RegisterUser struct {
	AuthBase
	DisplayName string
}

type LoginUser struct {
	AuthBase
}
