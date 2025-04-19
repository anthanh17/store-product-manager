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
	s.logger.Info("API call: createProduct")

	var req createProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		s.logger.Sugar().Infof("\nInvalid request data for createProduct: %v\n", err)
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
	product, err := s.store.CreateProductTx(ctx, arg)
	if err != nil {
		s.logger.Sugar().Infof("\nFailed to create product: %v\n", err)
		if db.ErrorCode(err) == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	s.logger.Sugar().Infof("\nProduct created successfully: %v\n", product.ProductId)

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Product created successfully",
	})
}

type getProductRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (s *Server) getProduct(ctx *gin.Context) {
	s.logger.Info("API call: getProduct")

	var req getProductRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		s.logger.Sugar().Infof("\nInvalid product ID: %v\n", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	s.logger.Sugar().Infof("\nFetching product detail product ID: %v\n", req.ID)
	product, err := s.store.GetProduct(ctx, req.ID)
	if err != nil {
		s.logger.Sugar().Infof("\nFailed to get product: %v\n", err)
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
		s.logger.Sugar().Infof("\nFailed to get reviews: %v\n", err)
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

	s.logger.Sugar().Infof("\nProduct details retrieved successfully: %v\n", req.ID)
	ctx.JSON(http.StatusOK, response)
}

type deleteProductRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (s *Server) deleteProduct(ctx *gin.Context) {
	s.logger.Info("API call: deleteProduct")

	var req deleteProductRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		s.logger.Sugar().Infof("\nInvalid product ID: %v\n", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Check products
	s.logger.Sugar().Infof("\nChecking if product exists, product ID: %v\n", req.ID)
	_, err := s.store.GetProduct(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Sugar().Infof("\nProduct not found for deletion, product ID: %v\n", req.ID)
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("product not found")))
			return
		}
		s.logger.Sugar().Infof("\nFailed to check product existence: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Delete product
	s.logger.Sugar().Infof("\nDeleting product, product ID: %v\n", req.ID)
	err = s.store.DeleteProduct(ctx, req.ID)
	if err != nil {
		s.logger.Sugar().Infof("\nFailed to delete product: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	s.logger.Sugar().Infof("\nProduct deleted successfully, product ID: %v\n", req.ID)
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
	s.logger.Info("API call: updateProduct")

	var reqURI struct {
		ID int32 `uri:"id" binding:"required,min=1"`
	}
	if err := ctx.ShouldBindUri(&reqURI); err != nil {
		s.logger.Sugar().Infof("\nInvalid product ID: %v\n", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get data from JSON body
	var req updateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		s.logger.Sugar().Infof("\nInvalid request data for updateProduct: %v\n", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// check product exits
	s.logger.Sugar().Infof("\nChecking if product exists, product ID: %v\n", reqURI.ID)
	_, err := s.store.GetProduct(ctx, reqURI.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Sugar().Infof("\nProduct not found for update, product ID: %v\n", reqURI.ID)
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("product not found")))
			return
		}
		s.logger.Sugar().Infof("\nFailed to check product existence: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Call transaction update product
	s.logger.Sugar().Infof("\nUpdating product, product ID: %v\n", reqURI.ID)
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
		s.logger.Sugar().Infof("\nFailed to update product: %v\n", err)
		if db.ErrorCode(err) == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	s.logger.Sugar().Infof("\nFetching updated product categories, product ID: %v\n", reqURI.ID)
	categories, err := s.store.GetProductCategories(ctx, reqURI.ID)
	if err != nil {
		s.logger.Sugar().Infof("\nFailed to get updated product categories: %v\n", err)
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

	s.logger.Sugar().Infof("\nProduct updated successfully, product ID: %v\n", reqURI.ID)
	ctx.JSON(http.StatusOK, response)
}

type listProductsRequest struct {
	Page              int32  `form:"page" binding:"required,min=1"`
	Limit             int32  `form:"limit" binding:"required,min=5,max=100"`
	Status            string `form:"status"`
	SearchProductName string `form:"search_product_name"`
}

func (s *Server) listProducts(ctx *gin.Context) {
	s.logger.Info("API call: listProducts")

	var req listProductsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		s.logger.Sugar().Infof("\nInvalid query parameters for listProducts: %v\n", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	offset := (req.Page - 1) * req.Limit
	s.logger.Sugar().Infof("\nFetching products with pagination, page: %v, limit: %v\n", req.Page, req.Limit)

	// Validate status if provided
	if req.Status != "" && req.Status != "IN_STOCK" && req.Status != "OUT_OF_STOCK" {
		s.logger.Sugar().Infof("\nInvalid status filter: %v\n", req.Status)
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid status, must be IN_STOCK or OUT_OF_STOCK")))
		return
	}

	// Log filter parameters if provided
	if req.Status != "" {
		s.logger.Sugar().Infof("\nFiltering by status: %v\n", req.Status)
	}
	if req.SearchProductName != "" {
		s.logger.Sugar().Infof("\nSearching for product name: %v\n", req.SearchProductName)
	}

	// Prepare search pattern for product name
	searchPattern := ""
	if req.SearchProductName != "" {
		searchPattern = "%" + req.SearchProductName + "%"
	}

	s.logger.Sugar().Infof("\nPreparing to query products with filters - Status: '%v', Search: '%v'\n", req.Status, req.SearchProductName)

	argListProducts := db.ListProductsWithFiltersParams{
		Limit:             req.Limit,
		Offset:            offset,
		Status:            req.Status,
		SearchProductName: searchPattern,
	}

	products, err := s.store.ListProductsWithFilters(ctx, argListProducts)
	if err != nil {
		s.logger.Sugar().Infof("\nFailed to list products: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	s.logger.Sugar().Infof("\nCounting total products with filters\n")
	countArg := db.CountProductsWithFiltersParams{
		Status:            req.Status,
		SearchProductName: searchPattern,
	}

	totalCount, err := s.store.CountProductsWithFilters(ctx, countArg)
	if err != nil {
		s.logger.Sugar().Infof("\nFailed to count products: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	s.logger.Sugar().Infof("\nProducts retrieved successfully, total count: %v\n", totalCount)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"products": products,
			"pagination": gin.H{
				"total": totalCount,
				"page":  req.Page,
				"limit": req.Limit,
				"pages": (totalCount + int64(req.Limit) - 1) / int64(req.Limit),
			},
		},
	})
}
