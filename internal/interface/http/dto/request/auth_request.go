package request

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=16"`
	Password string `json:"password" validate:"required,min=6,max=32"`
}

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Username       string `json:"username" validate:"required,min=3,max=16"`
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=6,max=32"`
	RepeatPassword string `json:"repeat_password" validate:"required,eqfield=Password"`
}
