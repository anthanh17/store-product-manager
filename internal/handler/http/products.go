package http

import (
	"database/sql"
	"fmt"
	"net/http"

	db "store-product-manager/internal/dataaccess/database/sqlc"

	"github.com/gin-gonic/gin"
)

type ProductStatus string

const (
	ProductStatusInStock    ProductStatus = "IN_STOCK"
	ProductStatusOutOfStock ProductStatus = "OUT_OF_STOCK"
)

// IsValid checks if the status is valid
func (s ProductStatus) IsValid() bool {
	return s == ProductStatusInStock || s == ProductStatusOutOfStock
}

type createProductRequest struct {
	Name          string        `json:"name" binding:"required"`
	Description   string        `json:"description"`
	Price         float64       `json:"price" binding:"required,min=0"`
	StockQuantity int32         `json:"stock_quantity" binding:"required,min=0"`
	Status        ProductStatus `json:"status" binding:"required,oneof=IN_STOCK OUT_OF_STOCK"`
	ImageURL      string        `json:"image_url"`
	CategoryIDs   []int32       `json:"category_ids"`
}

func (s *Server) createProduct(ctx *gin.Context) {
	var req createProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateProductTxParams{
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		StockQuantity: req.StockQuantity,
		Status:        string(req.Status),
		ImageURL:      req.ImageURL,
		CategoryIDs:   req.CategoryIDs,
	}

	// Call transaction create product
	_, err := s.store.CreateProductTx(ctx, arg)
	if err != nil {
		s.logger.Info("cannot CrecreateProductateUser")
		if db.ErrorCode(err) == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Product created successfully",
	})
}

type getProductRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (s *Server) getProduct(ctx *gin.Context) {
	var req getProductRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	product, err := s.store.GetProduct(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("product not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Get the product's review list
	reviews, err := s.store.GetProductReviews(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := gin.H{
		"status": "success",
		"data": gin.H{
			"id":             product.ID,
			"name":           product.Name,
			"description":    product.Description.String,
			"price":          product.Price,
			"stock_quantity": product.StockQuantity,
			"status":         product.Status,
			"image_url":      product.ImageUrl.String,
			"categories":     product.Categories,
			"reviews":        reviews,
			"created_at":     product.CreatedAt,
			"updated_at":     product.UpdatedAt,
		},
	}

	ctx.JSON(http.StatusOK, response)
}

type deleteProductRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (s *Server) deleteProduct(ctx *gin.Context) {
	var req deleteProductRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Check products
	_, err := s.store.GetProduct(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("product not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Delete product
	err = s.store.DeleteProduct(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Product deleted successfully",
	})
}

type updateProductRequest struct {
	Name          string  `json:"name" binding:"required"`
	Description   string  `json:"description"`
	Price         float64 `json:"price" binding:"required,gt=0"`
	StockQuantity int32   `json:"stock_quantity" binding:"required,gte=0"`
	Status        string  `json:"status" binding:"required,oneof=IN_STOCK OUT_OF_STOCK"`
	ImageUrl      string  `json:"image_url"`
	CategoryIds   []int32 `json:"category_ids"`
}

func (s *Server) updateProduct(ctx *gin.Context) {
	var reqURI struct {
		ID int32 `uri:"id" binding:"required,min=1"`
	}
	if err := ctx.ShouldBindUri(&reqURI); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get data from JSON body
	var req updateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// check product exits
	_, err := s.store.GetProduct(ctx, reqURI.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("product not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Call transaction update product
	arg := db.UpdateProductTxParams{
		ID:            reqURI.ID,
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		StockQuantity: req.StockQuantity,
		Status:        req.Status,
		ImageURL:      req.ImageUrl,
		CategoryIDs:   req.CategoryIds,
	}

	product, err := s.store.UpdateProductTx(ctx, arg)
	if err != nil {
		s.logger.Info("cannot update product")
		if db.ErrorCode(err) == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	categories, err := s.store.GetProductCategories(ctx, reqURI.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := gin.H{
		"status":  "success",
		"message": "Product updated successfully",
		"data": gin.H{
			"id":             product.ID,
			"name":           product.Name,
			"description":    product.Description,
			"price":          product.Price,
			"stock_quantity": product.StockQuantity,
			"status":         product.Status,
			"image_url":      product.ImageUrl,
			"categories":     categories,
			"updated_at":     product.UpdatedAt,
		},
	}

	ctx.JSON(http.StatusOK, response)
}
