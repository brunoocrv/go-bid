package users

import (
	"context"

	"github.com/brunoocrv/go-bid/internal/validator"
)

type SignInUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (req SignInUserReq) Valid(context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.Matches(req.Email, validator.EmailRX), "email", "email must be a valid email address")
	eval.CheckField(validator.NotBlank(req.Password), "password", "password must be provided")

	return eval
}
