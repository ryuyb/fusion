package command

type CreateUserCommand struct {
	Username string
	Email    string
	Password string
}

type UpdateUserCommand struct {
	*CreateUserCommand

	ID int64
}
