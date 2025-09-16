package repository

type User struct {
	Login    string
	Password string
}

type UserRepository struct {
	users map[string]User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}
