package http

import (
	"fmt"
	"sync"
	"time"

	"store-product-manager/configs"
	"store-product-manager/internal/dataaccess/cache"
	db "store-product-manager/internal/dataaccess/database/sqlc"
	"store-product-manager/internal/handler/token"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Server serves HTTP requests for our service.
type Server struct {
	config       configs.Config
	store        db.Store
	tokenMaker   token.Maker
	router       *gin.Engine
	sessionCache cache.SessionCache
	logger       *zap.Logger
	mu           *sync.Mutex
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(config configs.Config, store db.Store, cachier cache.Cachier, logger *zap.Logger) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.Token.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:       config,
		store:        store,
		tokenMaker:   tokenMaker,
		sessionCache: cache.NewSessionCache(cachier, logger),
		logger:       logger,
		mu:           new(sync.Mutex),
	}

	// Router
	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	// Use CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Authorization", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "accept", "origin", "Cache-Control", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           12 * time.Hour,
	}))

	router.POST("/api/auth/register", server.createUser)
	router.POST("/api/auth/login", server.loginUser)
	router.POST("/tokens/renew_access", server.renewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.GET("/ping", server.ping)
	authRoutes.GET("/top", server.top)
	authRoutes.GET("/count", server.count)

	// product routes
	authRoutes.GET("/api/products", server.listProducts)
	authRoutes.POST("/api/products", server.createProduct)
	authRoutes.GET("/api/products/:id", server.getProduct)
	authRoutes.PUT("/api/products/:id", server.updateProduct)
	authRoutes.DELETE("/api/products/:id", server.deleteProduct)

	// category routes
	authRoutes.POST("/api/categories", server.createCategory)
	authRoutes.GET("/api/categories", server.listCategories)
	authRoutes.GET("/api/categories/:id", server.getCategory)
	authRoutes.PUT("/api/categories/:id", server.updateCategory)
	authRoutes.DELETE("/api/categories/:id", server.deleteCategory)

	// dashboard api
	authRoutes.GET("/api/dashboard/summary", server.getDashboardSummary)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
