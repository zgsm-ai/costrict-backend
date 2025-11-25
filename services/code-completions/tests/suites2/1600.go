package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

// 简化版缓存系统实现

// CacheEntry 缓存条目
type CacheEntry struct {
	Key        string      `json:"key"`
	Value      interface{} `json:"value"`
	Expiration int64       `json:"expiration"`  // timestamp
	AccessedAt int64       `json:"accessed_at"` // timestamp
	Size       int64       `json:"size"`        // in bytes
}

// CacheStats 缓存统计
type CacheStats struct {
	TotalEntries int64   `json:"total_entries"`
	TotalSize    int64   `json:"total_size"`
	TotalHits    int64   `json:"total_hits"`
	TotalMisses  int64   `json:"total_misses"`
	HitRate      float64 `json:"hit_rate"`
}

// SimpleCache 简单缓存
type SimpleCache struct {
	entries map[string]*CacheEntry
	stats   CacheStats
	config  CacheConfig
	mutex   sync.RWMutex
}

// CacheConfig 缓存配置
type CacheConfig struct {
	DefaultExpiration int64  `json:"default_expiration"` // in seconds
	MaxMemory         int64  `json:"max_memory"`         // in bytes
	EvictionPolicy    string `json:"eviction_policy"`    // "lru", "random"
}

// NewSimpleCache 创建新的简单缓存
func NewSimpleCache(config CacheConfig) *SimpleCache {
	return &SimpleCache{
		entries: make(map[string]*CacheEntry),
		stats:   CacheStats{},
		config:  config,
	}
}

// Set 设置缓存
func (sc *SimpleCache) Set(key string, value interface{}, expiration int64) error {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	// 检查是否需要淘汰
	entrySize := int64(16) // 简化处理，假设每个条目固定大小
	if sc.stats.TotalSize+entrySize > sc.config.MaxMemory {
		sc.evictItems(entrySize)
	}

	// 创建缓存条目
	entry := &CacheEntry{
		Key:        key,
		Value:      value,
		Expiration: time.Now().Unix() + expiration,
		AccessedAt: time.Now().Unix(),
		Size:       entrySize,
	}

	// 检查是否已存在
	if _, exists := sc.entries[key]; exists {
		sc.stats.TotalSize -= entrySize
	} else {
		sc.stats.TotalEntries++
	}

	sc.entries[key] = entry
	sc.stats.TotalSize += entrySize

	return nil
}

// Get 获取缓存
func (sc *SimpleCache) Get(key string) (interface{}, bool) {
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()

	entry, exists := sc.entries[key]
	if !exists {
		sc.stats.TotalMisses++
		return nil, false
	}

	// 检查是否过期
	if entry.Expiration < time.Now().Unix() {
		sc.stats.TotalMisses++
		go func() {
			sc.mutex.Lock()
			delete(sc.entries, key)
			sc.stats.TotalEntries--
			sc.stats.TotalSize -= entry.Size
			sc.mutex.Unlock()
		}()
		return nil, false
	}

	// 更新访问时间和命中次数
	entry.AccessedAt = time.Now().Unix()
	sc.stats.TotalHits++

	// 更新命中率
	if sc.stats.TotalHits+sc.stats.TotalMisses > 0 {
		sc.stats.HitRate = float64(sc.stats.TotalHits) / float64(sc.stats.TotalHits+sc.stats.TotalMisses)
	}

	return entry.Value, true
}

// Delete 删除缓存
func (sc *SimpleCache) Delete(key string) bool {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	entry, exists := sc.entries[key]
	if !exists {
		return false
	}

	delete(sc.entries, key)
	sc.stats.TotalEntries--
	sc.stats.TotalSize -= entry.Size

	return true
}

// Clear 清空缓存
func (sc *SimpleCache) Clear() {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	sc.entries = make(map[string]*CacheEntry)
	sc.stats = CacheStats{}
}

// GetStats 获取统计信息
func (sc *SimpleCache) GetStats() CacheStats {
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()

	// 计算命中率
	if sc.stats.TotalHits+sc.stats.TotalMisses > 0 {
		sc.stats.HitRate = float64(sc.stats.TotalHits) / float64(sc.stats.TotalHits+sc.stats.TotalMisses)
	}

	return sc.stats
}

// evictItems 淘汰条目
func (sc *SimpleCache) evictItems(requiredSpace int64) {
	// 根据淘汰策略淘汰条目
	switch sc.config.EvictionPolicy {
	case "lru":
		sc.evictLRU(requiredSpace)
	default: // "random"
		sc.evictRandom(requiredSpace)
	}
}

// evictLRU 使用LRU算法淘汰条目
func (sc *SimpleCache) evictLRU(requiredSpace int64) {
	freedSpace := int64(0)

	// 创建键列表
	keys := make([]string, 0, len(sc.entries))
	for key := range sc.entries {
		keys = append(keys, key)
	}

	// 按访问时间排序
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if sc.entries[keys[i]].AccessedAt > sc.entries[keys[j]].AccessedAt {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}

	// 淘汰条目直到释放足够空间
	for _, key := range keys {
		if freedSpace >= requiredSpace {
			break
		}

		entry := sc.entries[key]
		delete(sc.entries, key)
		sc.stats.TotalEntries--
		sc.stats.TotalSize -= entry.Size
		freedSpace += entry.Size
	}
}

// evictRandom 随机淘汰条目
func (sc *SimpleCache) evictRandom(requiredSpace int64) {
	freedSpace := int64(0)

	// 创建键列表
	keys := make([]string, 0, len(sc.entries))
	for key := range sc.entries {
		keys = append(keys, key)
	}

	// 随机淘汰条目直到释放足够空间
	for freedSpace < requiredSpace && len(keys) > 0 {
		// 随机选择一个键
		idx := rand.Intn(len(keys))
		key := keys[idx]

		// 淘汰条目
		entry := sc.entries[key]
		delete(sc.entries, key)
		sc.stats.TotalEntries--
		sc.stats.TotalSize -= entry.Size
		freedSpace += entry.Size

		// 从列表中移除
		keys = append(keys[:idx], keys[idx+1:]...)
	}
}

func main() {
	// 创建缓存配置
	config := CacheConfig{
		DefaultExpiration: 3600,             // 1小时
		MaxMemory:         10 * 1024 * 1024, // 10MB
		EvictionPolicy:    "lru",
	}
	<｜fim▁hole｜>
	// 创建简单缓存
	cache := NewSimpleCache(config)

	// 设置一些测试数据
	for i := 1; i <= 100; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := map[string]interface{}{
			"id":    i,
			"name":  fmt.Sprintf("Item %d", i),
			"value": rand.Intn(1000),
		}

		if err := cache.Set(key, value, config.DefaultExpiration); err != nil {
			log.Printf("设置缓存失败: %v", err)
		}
	}

	// 获取一些测试数据
	for i := 1; i <= 10; i++ {
		key := fmt.Sprintf("key-%d", i)
		if value, exists := cache.Get(key); exists {
			log.Printf("获取缓存成功: %s = %v", key, value)
		} else {
			log.Printf("获取缓存失败: %s", key)
		}
	}

	// 获取统计信息
	stats := cache.GetStats()
	log.Printf("缓存统计: %+v", stats)

	// 等待一段时间
	time.Sleep(1 * time.Second)

	fmt.Println("简化版缓存系统示例完成")
}
