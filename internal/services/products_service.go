package services

import (
	"context"
	"errors"
	"time"

	"github.com/brunoocrv/go-bid/internal/store/pgstore"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductsService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

func NewProductsService(pool *pgxpool.Pool) ProductsService {
	return ProductsService{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

var ErrProductNotFound = errors.New("product not found")

func (ps ProductsService) CreateProduct(
	ctx context.Context,
	sellerId uuid.UUID, name, description string, basePrice float64, auctionEnd time.Time,
) (uuid.UUID, error) {
	id, err := ps.queries.CreateProduct(ctx, pgstore.CreateProductParams{
		SellerID:   sellerId,
		Name:       name,
		BasePrice:  basePrice,
		AuctionEnd: auctionEnd,
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

func (ps ProductsService) GetProductById(ctx context.Context, id uuid.UUID) (pgstore.GetProductByIdRow, error) {
	product, err := ps.queries.GetProductById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgstore.GetProductByIdRow{}, ErrProductNotFound
		}

		return pgstore.GetProductByIdRow{}, err
	}

	return product, nil
}
