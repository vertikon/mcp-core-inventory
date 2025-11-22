package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// StockCache implements stock availability caching using Redis
type StockCache struct {
	client *redis.Client
}

// NewStockCache creates a new Redis stock cache
func NewStockCache(client *redis.Client) *StockCache {
	return &StockCache{client: client}
}

// GetAvailable gets available stock from cache
func (c *StockCache) GetAvailable(ctx context.Context, sku, location string) (int64, bool) {
	key := cacheKey(sku, location)
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return 0, false
	}

	quantity, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, false
	}

	return quantity, true
}

// SetAvailable sets available stock in cache
func (c *StockCache) SetAvailable(ctx context.Context, sku, location string, quantity int64, ttlSeconds int) error {
	key := cacheKey(sku, location)
	return c.client.Set(ctx, key, quantity, time.Duration(ttlSeconds)*time.Second).Err()
}

// Invalidate invalidates the cache for a SKU/location
func (c *StockCache) Invalidate(ctx context.Context, sku, location string) error {
	key := cacheKey(sku, location)
	return c.client.Del(ctx, key).Err()
}

// InvalidateByPattern invalidates cache entries matching a pattern
func (c *StockCache) InvalidateByPattern(ctx context.Context, pattern string) error {
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return c.client.Del(ctx, keys...).Err()
	}

	return nil
}

// GetStockInfo gets full stock info from cache (if available)
func (c *StockCache) GetStockInfo(ctx context.Context, sku, location string) (*StockInfo, bool) {
	key := infoCacheKey(sku, location)
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, false
	}

	var info StockInfo
	if err := json.Unmarshal([]byte(val), &info); err != nil {
		return nil, false
	}

	return &info, true
}

// SetStockInfo sets full stock info in cache
func (c *StockCache) SetStockInfo(ctx context.Context, sku, location string, info *StockInfo, ttlSeconds int) error {
	key := infoCacheKey(sku, location)
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, time.Duration(ttlSeconds)*time.Second).Err()
}

// StockInfo represents cached stock information
type StockInfo struct {
	SKU       string `json:"sku"`
	Location  string `json:"location"`
	Available int64  `json:"available"`
	Reserved  int64  `json:"reserved"`
	Total     int64  `json:"total"`
}

func cacheKey(sku, location string) string {
	return fmt.Sprintf("stock:available:%s:%s", sku, location)
}

func infoCacheKey(sku, location string) string {
	return fmt.Sprintf("stock:info:%s:%s", sku, location)
}
