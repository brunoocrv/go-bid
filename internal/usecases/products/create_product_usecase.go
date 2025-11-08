package products

import (
	"context"
	"time"

	"github.com/brunoocrv/go-bid/internal/validator"
)

type CreateProductReq struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	BasePrice   float64   `json:"base_price"`
	AuctionEnd  time.Time `json:"auction_end"`
}

const minAuctionDuration = 2 * time.Hour

func (req CreateProductReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.NotBlank(req.Name), "name", "this field cannot be blank")
	eval.CheckField(validator.NotBlank(req.Description), "name", "this field cannot be blank")
	eval.CheckField(
		validator.MinChars(req.Description, 10) && validator.MaxChars(req.Description, 255),
		"description",
		"this field must be between 10 and 255 characters")
	eval.CheckField(req.BasePrice > 0, "base_price", "this field must be greater than 0")
	eval.CheckField(req.AuctionEnd.Sub(time.Now()) >= minAuctionDuration, "auction_end", "this field must be greater than or equal to the minimum auction duration")

	return eval
}
