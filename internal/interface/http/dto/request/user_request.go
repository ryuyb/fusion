package request

type CreateUserRequest struct {
	Username       string `json:"username" validate:"required,min=3,max=16"`
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=6,max=32"`
	RepeatPassword string `json:"repeat_password" validate:"required,eqfield=Password"`
	Status         string `json:"status" validate:"required,oneof=active inactive banned"`
}

type UpdateUserRequest struct {
	ID int64 `json:"id" validate:"required,gte=1"`
	CreateUserRequest
}
