package cache

import (
	"context"
	"fmt"
	"time"

	"store-product-manager/internal/utils"
	"go.uber.org/zap"
)

type SessionCache interface {
	GetSession(ctx context.Context, id string) (SessionType, error)
	SetSession(ctx context.Context, id string, data SessionType) error

	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, data any) error
	// SetNX
	SetPingLock(ctx context.Context, key string, data any) (bool, error)
	Del(ctx context.Context, key string) error

	// Add on
	CheckRateLimit(ctx context.Context, key string, limit int, duration time.Duration) (bool, error)
	// Increase the number of calls in Sorted Set
	IncreaseTopCalls(ctx context.Context, key string, username string) error
	// Get the top 10 users with the most calls
	Top10UsersCalling(ctx context.Context, key string) ([]string, error)
	// Use hyperloglog to store the approximate number of api /ping callers
	AddUsernameHyperLogLog(ctx context.Context, key string, username string) error
	CountUsernameHyperLogLog(ctx context.Context, key string) (int64, error)
}

type SessionType struct {
	SessionID string `json:"sessionId"`
	Username  string `json:"username"`
}

type sessionCache struct {
	cachier Cachier
	logger  *zap.Logger
}

func NewSessionCache(cachier Cachier, logger *zap.Logger) SessionCache {
	return &sessionCache{
		cachier: cachier,
		logger:  logger,
	}
}

func (s sessionCache) getSessionCacheKey(id string) string {
	return fmt.Sprintf("session_cache_key:%s", id)
}

func (s sessionCache) GetSession(ctx context.Context, id string) (SessionType, error) {
	logger := utils.LoggerWithContext(ctx, s.logger).With(zap.String("id", id))

	// get key cache
	cacheKey := s.getSessionCacheKey(id)

	// Get cache data
	cacheEntry, err := s.cachier.Get(ctx, cacheKey)
	if err != nil {
		// logger.With(zap.Error(err)).Error("failed to get session key cache")
		logger.Info("failed to get session key cache")
		return SessionType{}, err
	}

	// If miss cache
	if cacheEntry == nil {
		return SessionType{}, ErrCacheMiss
	}

	// Check data type session
	sessionData, ok := cacheEntry.(SessionType)
	if !ok {
		logger.Error("cache entry is not of SessionType")
		return SessionType{}, nil
	}

	return sessionData, nil
}

func (s sessionCache) Get(ctx context.Context, key string) (any, error) {
	logger := utils.LoggerWithContext(ctx, s.logger).With(zap.String("key", key))

	// Get cache data
	cacheEntry, err := s.cachier.Get(ctx, key)
	if err != nil {
		logger.Info("failed to get session key cache")
		return "", err
	}

	// If miss cache
	if cacheEntry == nil {
		logger.Info("miss cache")
		return "", ErrCacheMiss
	}

	return cacheEntry, nil
}

func (s sessionCache) SetSession(ctx context.Context, id string, data SessionType) error {
	logger := utils.LoggerWithContext(ctx, s.logger).With(zap.String("id", id))

	// get key cache
	cacheKey := s.getSessionCacheKey(id)

	if err := s.cachier.Set(ctx, cacheKey, data, 0); err != nil {
		// logger.With(zap.Error(err)).Error("failed to insert token public key into cache")
		logger.Info("failed to insert token public key into cache")
		return err
	}

	return nil
}

func (s sessionCache) Set(ctx context.Context, key string, data any) error {
	logger := utils.LoggerWithContext(ctx, s.logger).With(zap.String("key", key))

	if err := s.cachier.Set(ctx, key, data, 0); err != nil {
		logger.Info("failed to insert token public key into cache")
		return err
	}

	return nil
}

func (s sessionCache) Del(ctx context.Context, key string) error {
	err := s.cachier.Del(ctx, key)
	if err != nil {
		s.logger.Info("failed to delete session key cache")
		return err
	}

	return nil
}

func (s sessionCache) SetPingLock(ctx context.Context, key string, data any) (bool, error) {
	ok, err := s.cachier.SetNX(ctx, key, data, 0)
	if err != nil {
		s.logger.Info("failed SetPingLock")
		return false, err
	}

	return ok, nil
}

func (s sessionCache) CheckRateLimit(ctx context.Context, key string, limit int, duration time.Duration) (bool, error) {
	cachier, ok := s.cachier.(*redisClient)
	if !ok {
		s.logger.Info("cachier is not redis")
		return false, fmt.Errorf("cachier is not redis")
	}

	// Check rate limit
	// Get all element in list [0:-1]
	val, err := cachier.redisClient.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		s.logger.Info("erro get list redis: " + err.Error())
		return false, fmt.Errorf("error LRange")
	}

	// Exits lenght list > limit
	if len(val) >= limit {
		return false, nil
	}

	// Add current timestamp to list
	err = cachier.redisClient.LPush(ctx, key, time.Now().Unix()).Err()
	if err != nil {
		s.logger.Info("error add list redis: " + err.Error())
		return true, fmt.Errorf("error add LPush list redis")
	}

	// Setting TTL for key (60s)
	err = cachier.redisClient.Expire(ctx, key, duration).Err()
	if err != nil {
		s.logger.Info("error Setting TTL for key: " + err.Error())
		return false, fmt.Errorf("error Setting TTL for key")
	}

	return true, nil
}

func (s sessionCache) IncreaseTopCalls(ctx context.Context, key string, username string) error {
	cachier, ok := s.cachier.(*redisClient)
	if !ok {
		s.logger.Info("cachier is not redis")
		return fmt.Errorf("cachier is not redis")
	}

	// Increase the number of calls in Sorted Set
	err := cachier.redisClient.ZIncrBy(ctx, key, 1, username).Err()
	if err != nil {
		s.logger.Info("error ZIncrBy:" + err.Error())
		return fmt.Errorf("error ZIncrBy")
	}

	return nil
}

func (s sessionCache) Top10UsersCalling(ctx context.Context, key string) ([]string, error) {
	cachier, ok := s.cachier.(*redisClient)
	if !ok {
		s.logger.Info("cachier is not redis")
		return nil, fmt.Errorf("cachier is not redis")
	}

	// Get the top 10 users with the most calls
	result, err := cachier.redisClient.ZRevRangeWithScores(ctx, "top_users", 0, 9).Result()
	if err != nil {
		s.logger.Info("error ZRevRangeWithScores:" + err.Error())
		return nil, fmt.Errorf("error ZRevRangeWithScores")
	}

	// Convert result -> username
	var topUsers []string
	for _, item := range result {
		topUsers = append(topUsers, item.Member.(string))
	}

	return topUsers, nil
}

func (s sessionCache) AddUsernameHyperLogLog(ctx context.Context, key string, username string) error {
	cachier, ok := s.cachier.(*redisClient)
	if !ok {
		s.logger.Info("cachier is not redis")
		return fmt.Errorf("cachier is not redis")
	}

	// Add username to HyperLogLog
	err := cachier.redisClient.PFAdd(ctx, key, username).Err()
	if err != nil {
		s.logger.Info("error PFAdd:" + err.Error())
		return fmt.Errorf("error PFAdd")
	}

	return nil
}

func (s sessionCache) CountUsernameHyperLogLog(ctx context.Context, key string) (int64, error) {
	cachier, ok := s.cachier.(*redisClient)
	if !ok {
		s.logger.Info("cachier is not redis")
		return -1, fmt.Errorf("cachier is not redis")
	}

	// Count Username in HyperLogLog
	count, err := cachier.redisClient.PFCount(ctx, key).Result()
	if err != nil {
		s.logger.Info("error PFCount:" + err.Error())
		return -1, fmt.Errorf("error PFCount")
	}

	return count, nil
}
