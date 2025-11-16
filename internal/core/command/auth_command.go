package command

type LoginCommand struct {
	Username string
	Password string
}

type RegisterCommand struct {
	Username string
	Email    string
	Password string
}
