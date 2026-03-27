package v1

type User struct {
	Id          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}

type UserResponse struct {
	Data User `json:"data"`
}
