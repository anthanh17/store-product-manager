package http

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"store-product-manager/internal/handler/token"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// AuthMiddleware creates a gin middleware for authorization
func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

// RateLimitMiddleware creates a gin middleware for rate limiting
func (s *Server) rateLimitMiddleware(limit int, duration time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get data by access token
		accessPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

		// Check rate limit
		rateLimitKey := fmt.Sprintf("rate_limit:%s:%s", ctx.FullPath(), accessPayload.Username)

		// Check if rate limit is exceeded
		ok, err := s.sessionCache.CheckRateLimit(ctx, rateLimitKey, limit, duration)
		if err != nil || !ok {
			ctx.JSON(
				http.StatusTooManyRequests,
				gin.H{"error": fmt.Sprintf("Rate limit exceeded: maximum %d requests in %v", limit, duration)},
			)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
