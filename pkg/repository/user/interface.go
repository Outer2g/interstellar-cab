package repository

type UserRepository interface {
	GetUser(email string) *User
	AddUser(email string, passwordhash string, isVip bool) bool
}
