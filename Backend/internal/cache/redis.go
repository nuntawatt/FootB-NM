package cache

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

// เขื่อมต่อ Redis ถ้าต่อไม่ติดจะไม่แครช แต่จะ log เตือน
func New(redisAddr string) *Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("warning: redis connection failed (%v); running in cache-bypass mode\n", err)
		return &Cache{client: nil}
	}

	log.Println("redis connected successfully")
	return &Cache{client: rdb}
}

// ดึงข้อมูลไจาก Cache และแปลง json ลงใน target จะคืน true ถ้าเจอข้อมูล
func (c *Cache) Get(ctx context.Context, key string, target any) bool {
	if c == nil || c.client == nil {
		return false
	}

	val, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Printf("cache get error for key %s\n", key, err)
		}
		return false
	}

	if err := json.Unmarshal(val, target); err != nil {
		return false
	}
	return true
}

// บันทึกข้อมูลลง Cache ในรูปแบบ json และกำหนด TTL
func (c *Cache) Set(ctx context.Context, key string, value any, ttl time.Duration) {
	if c == nil || c.client == nil {
		return
	}

	bytes, err := json.Marshal(value)
	if err != nil {
		return
	}

	if err := c.client.Set(ctx, key, bytes, ttl).Err(); err != nil {
		log.Printf("cache set error for key %s\n", key, err)
	}
}

// ลบแคชทั้งหมดที่ชินต้นด้วย prefix
func (c *Cache) DeletePrefix(ctx context.Context, prefix string) {
	if c == nil || c.client == nil {
		return
	}

	iter := c.client.Scan(ctx, 0, prefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		_ = c.client.Del(ctx, iter.Val()).Err()
	}
}
