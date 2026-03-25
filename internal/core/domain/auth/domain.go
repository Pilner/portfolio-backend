package domain

type User struct {
	Id          string
	Email       string
	DisplayName string
}

type AddUser struct {
	Email       string
	Password    string
	DisplayName string
}
