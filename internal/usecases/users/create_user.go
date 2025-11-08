package users

import (
	"context"

	"github.com/brunoocrv/go-bid/internal/validator"
)

type CreateUserReq struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

func (req CreateUserReq) Valid(context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.Matches(req.Email, validator.EmailRX), "email", "this field must be a valid email address")
	eval.CheckField(validator.NotBlank(req.Email), "email", "this field cannot be empty")
	eval.CheckField(validator.NotBlank(req.UserName), "user_name", "this field cannot be empty")

	eval.CheckField(validator.NotBlank(req.Bio), "bio", "this field cannot be empty")
	eval.CheckField(
		validator.MinChars(req.Bio, 10) && validator.MaxChars(req.Bio, 255),
		"bio",
		"this field must be at least 10 characters long and not exceed 255 characters")

	eval.CheckField(validator.MinChars(req.Password, 8), "password", "this field must be at least 8 characters long")

	return eval
}
