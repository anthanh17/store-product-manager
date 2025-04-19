package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// CREATE
type CreateProductTxParams struct {
	Name          string
	Description   string
	Price         float64
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
			Price:         arg.Price,
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

// UpdateProductTxParams chứa các tham số cần thiết để cập nhật sản phẩm
type UpdateProductTxParams struct {
	ID            int32   `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	StockQuantity int32   `json:"stock_quantity"`
	Status        string  `json:"status"`
	ImageURL      string  `json:"image_url"`
	CategoryIDs   []int32 `json:"category_ids"`
}

// UpdateProductTxResult chứa kết quả của transaction cập nhật sản phẩm
type UpdateProductTxResult struct {
	ID            int32     `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Price         float64   `json:"price"`
	StockQuantity int32     `json:"stock_quantity"`
	Status        string    `json:"status"`
	ImageUrl      string    `json:"image_url"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (store *SQLStore) UpdateProductTx(ctx context.Context, arg UpdateProductTxParams) (UpdateProductTxResult, error) {
	var result UpdateProductTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// 1. update product
		updateArg := UpdateProductParams{
			ID:            arg.ID,
			Name:          arg.Name,
			Description:   pgtype.Text{String: arg.Description, Valid: true},
			Price:         arg.Price,
			StockQuantity: arg.StockQuantity,
			Status:        arg.Status,
			ImageUrl:      pgtype.Text{String: arg.ImageURL, Valid: arg.ImageURL != ""},
		}

		product, err := q.UpdateProduct(ctx, updateArg)
		if err != nil {
			return err
		}

		// save result
		result.ID = product.ID
		result.Name = product.Name
		result.Description = product.Description.String
		result.Price = product.Price
		result.StockQuantity = product.StockQuantity
		result.Status = product.Status
		result.ImageUrl = product.ImageUrl.String
		result.UpdatedAt = product.UpdatedAt

		// 2. delete all category
		err = q.DeleteProductCategories(ctx, arg.ID)
		if err != nil {
			return err
		}

		// 3. add new category
		for _, categoryID := range arg.CategoryIDs {
			err = q.AddProductCategory(ctx, AddProductCategoryParams{
				ProductID:  arg.ID,
				CategoryID: categoryID,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	return result, err
}
