package db

import (
	"context"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
)

// CREATE
type CreateProductTxParams struct {
	Name          string
	Description   string
	Price         int32
	StockQuantity int32
	Status        string
	ImageURL      string
	CategoryIDs   []int32
}

type CreateProductTxResult struct {
	ProductId int32
}

func (store *SQLStore) CreateProductTx(ctx context.Context, arg CreateProductTxParams) (CreateProductTxResult, error) {
	var result CreateProductTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// table: products
		argProduct := CreateProductParams{
			Name:          arg.Name,
			Description:   pgtype.Text{String: arg.Description, Valid: true},
			Price:         pgtype.Numeric{Int: big.NewInt(int64(arg.Price)), Valid: true},
			StockQuantity: arg.StockQuantity,
			Status:        arg.Status,
			ImageUrl:      pgtype.Text{String: arg.ImageURL, Valid: true},
		}
		product, err := q.CreateProduct(ctx, argProduct)
		if err != nil {
			return err
		}
		result.ProductId = product.ID

		// table: product_categories
		for _, categoryID := range arg.CategoryIDs {
			argProductCategory := CreateProductCategoryParams{
				ProductID:  result.ProductId,
				CategoryID: categoryID,
			}
			_, err = q.CreateProductCategory(ctx, argProductCategory)
			if err != nil {
				return err
			}
		}

		return nil

	})
	return result, err
}
