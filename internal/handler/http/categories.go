package http

import (
	"database/sql"
	"fmt"
	"net/http"

	db "store-product-manager/internal/dataaccess/database/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type createCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (s *Server) createCategory(ctx *gin.Context) {
	s.logger.Info("API call: createCategory")

	var req createCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		s.logger.Sugar().Infof("\nInvalid request data for createCategory: %v\n", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateCategoryParams{
		Name:        req.Name,
		Description: pgtype.Text{String: req.Description, Valid: true},
	}

	category, err := s.store.CreateCategory(ctx, arg)
	if err != nil {
		s.logger.Sugar().Infof("\nFailed to create category: %v\n", err)
		if db.ErrorCode(err) == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	s.logger.Sugar().Infof("\nCategory created successfully, category ID: %v\n", category.ID)
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Category created successfully",
	})
}

type listCategoriesRequest struct {
	Page  int32 `form:"page" binding:"required,min=1"`
	Limit int32 `form:"limit" binding:"required,min=5,max=100"`
}

func (s *Server) listCategories(ctx *gin.Context) {
	s.logger.Info("API call: listCategories")

	var req listCategoriesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		s.logger.Sugar().Infof("\nInvalid query parameters for listCategories: %v\n", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	offset := (req.Page - 1) * req.Limit
	s.logger.Sugar().Infof("\nFetching categories with pagination, page: %v, limit: %v\n", req.Page, req.Limit)

	argListCategory := db.ListCategoriesWithPaginationParams{
		Limit:  req.Limit,
		Offset: offset,
	}
	categories, err := s.store.ListCategoriesWithPagination(ctx, argListCategory)
	if err != nil {
		s.logger.Sugar().Infof("\nFailed to list categories: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	s.logger.Sugar().Infof("\nCounting total categories\n")
	totalCount, err := s.store.CountCategories(ctx)
	if err != nil {
		s.logger.Sugar().Infof("\nFailed to count categories: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	s.logger.Sugar().Infof("\nCategories retrieved successfully, total count: %v\n", totalCount)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"categories": categories,
			"pagination": gin.H{
				"total": totalCount,
				"page":  req.Page,
				"limit": req.Limit,
				"pages": (totalCount + int64(req.Limit) - 1) / int64(req.Limit),
			},
		},
	})
}

type getCategoryRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (s *Server) getCategory(ctx *gin.Context) {
	s.logger.Info("API call: getCategory")

	var req getCategoryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		s.logger.Sugar().Infof("\nInvalid category ID: %v\n", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	s.logger.Sugar().Infof("\nFetching category details, category ID: %v\n", req.ID)
	category, err := s.store.GetCategory(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Sugar().Infof("\nCategory not found, category ID: %v\n", req.ID)
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("category not found")))
			return
		}
		s.logger.Sugar().Infof("\nFailed to get category: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	s.logger.Sugar().Infof("\nCategory details retrieved successfully, category ID: %v\n", req.ID)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   category,
	})
}

type updateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (s *Server) updateCategory(ctx *gin.Context) {
	s.logger.Info("API call: updateCategory")

	var reqURI getCategoryRequest
	if err := ctx.ShouldBindUri(&reqURI); err != nil {
		s.logger.Sugar().Infof("\nInvalid category ID: %v\n", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var reqJSON updateCategoryRequest
	if err := ctx.ShouldBindJSON(&reqJSON); err != nil {
		s.logger.Sugar().Infof("\nInvalid request data for updateCategory: %v\n", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	s.logger.Sugar().Infof("\nUpdating category, category ID: %v\n", reqURI.ID)
	arg := db.UpdateCategoryParams{
		ID:          reqURI.ID,
		Name:        reqJSON.Name,
		Description: pgtype.Text{String: reqJSON.Description, Valid: true},
	}

	category, err := s.store.UpdateCategory(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Sugar().Infof("\nCategory not found for update, category ID: %v\n", reqURI.ID)
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("category not found")))
			return
		}
		if db.ErrorCode(err) == db.UniqueViolation {
			s.logger.Sugar().Infof("\nCategory name already exists: %v\n", reqJSON.Name)
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		s.logger.Sugar().Infof("\nFailed to update category: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	s.logger.Sugar().Infof("\nCategory updated successfully, category ID: %v\n", reqURI.ID)
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Category updated successfully",
		"data":    category,
	})
}

type deleteCategoryRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (s *Server) deleteCategory(ctx *gin.Context) {
	s.logger.Info("API call: deleteCategory")

	var req deleteCategoryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		s.logger.Sugar().Infof("\nInvalid category ID: %v\n", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Check if category exists
	s.logger.Sugar().Infof("\nChecking if category exists, category ID: %v\n", req.ID)
	_, err := s.store.GetCategory(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Sugar().Infof("\nCategory not found for deletion, category ID: %v\n", req.ID)
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("category not found")))
			return
		}
		s.logger.Sugar().Infof("\nFailed to check category existence: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	s.logger.Sugar().Infof("\nDeleting category, category ID: %v\n", req.ID)
	err = s.store.DeleteCategory(ctx, req.ID)
	if err != nil {
		s.logger.Sugar().Infof("\nFailed to delete category: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	s.logger.Sugar().Infof("\nCategory deleted successfully, category ID: %v\n", req.ID)
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Category deleted successfully",
	})
}
