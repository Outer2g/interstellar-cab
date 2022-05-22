package repository

type User struct {
	Email        string
	Passwordhash string
	Vip          bool
}

// checks if the password hash is valid
func (u *User) ValidatePasswordHash(pswdhash string) bool {
	return u.Passwordhash == pswdhash
}

type database struct {
	users map[string]User
}

func NewUserInMemoryDatabase() *database {
	return &database{map[string]User{}}
}

func (repo *database) GetUser(email string) *User {
	//needs to be replaces using Database
	if user, present := repo.users[email]; present {
		return &user
	}
	return nil
}

// returns true if the user already exists, false otherwise
func (repo *database) AddUser(email string, passwordhash string, isVip bool) bool {
	if _, present := repo.users[email]; present {
		return true
	}
	repo.users[email] = User{email, passwordhash, isVip}
	return false
}
