package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategorySummary struct {
	ID           int32  `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	ProductCount int64  `json:"product_count"`
}

type StatusSummary struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type DashboardSummary struct {
	TotalCategories int64             `json:"total_categories"`
	TotalProducts   int64             `json:"total_products"`
	Categories      []CategorySummary `json:"categories"`
	StatusSummary   []StatusSummary   `json:"status_summary"`
}

func (s *Server) getDashboardSummary(ctx *gin.Context) {
	s.logger.Info("API call: getDashboardSummary")

	// Get category summary
	s.logger.Sugar().Infof("\nFetching category summary\n")
	categorySummary, err := s.store.GetCategorySummary(ctx)
	if err != nil {
		s.logger.Sugar().Infof("\nFailed to get category summary: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Get total categories
	s.logger.Sugar().Infof("\nFetching total categories\n")
	totalCategories, err := s.store.GetTotalCategories(ctx)
	if err != nil {
		s.logger.Sugar().Infof("\nFailed to get total categories: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Get total products
	s.logger.Sugar().Infof("\nFetching total products\n")
	totalProducts, err := s.store.GetTotalProducts(ctx)
	if err != nil {
		s.logger.Sugar().Infof("\nFailed to get total products: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Get product status summary
	s.logger.Sugar().Infof("\nFetching product status summary\n")
	statusSummary, err := s.store.GetProductStatusSummary(ctx)
	if err != nil {
		s.logger.Sugar().Infof("\nFailed to get product status summary: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Convert to response format
	categories := make([]CategorySummary, len(categorySummary))
	for i, cat := range categorySummary {
		categories[i] = CategorySummary{
			ID:           cat.ID,
			Name:         cat.Name,
			Description:  cat.Description.String,
			ProductCount: cat.ProductCount,
		}
	}

	statuses := make([]StatusSummary, len(statusSummary))
	for i, status := range statusSummary {
		statuses[i] = StatusSummary{
			Status: status.Status,
			Count:  status.Count,
		}
	}

	response := DashboardSummary{
		TotalCategories: totalCategories,
		TotalProducts:   totalProducts,
		Categories:      categories,
		StatusSummary:   statuses,
	}

	s.logger.Sugar().Infof("\nDashboard summary retrieved successfully\n")
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   response,
	})
}
