package cache

import (
	"errors"
)

// Cache описывает интерфейс для кеширования.
type Cache interface {
	Get(key string) (string, error)
	Set(key, value string, ttl int) error
	Delete(key string) error
}

// MultiLevelCache представляет многослойный кеш.
type MultiLevelCache struct {
	caches []Cache
}

// NewMultiLevelCache создает новый многослойный кеш.
func NewMultiLevelCache(caches ...Cache) *MultiLevelCache {
	return &MultiLevelCache{
		caches: caches,
	}
}

// Get получает значение из кеша. Поиск идет сверху вниз.
func (m *MultiLevelCache) Get(key string) (string, error) {
	for _, cache := range m.caches {
		if value, err := cache.Get(key); err == nil && value != "" {
			// Если найдено, обновляем более высокие уровни кеша
			_ = m.Set(key, value, 300)
			return value, nil
		}
	}
	return "", errors.New("key not found")
}

// Set сохраняет значение во все уровни кеша.
func (m *MultiLevelCache) Set(key string, value string, ttl int) error {
	for _, cache := range m.caches {
		if err := cache.Set(key, value, ttl); err != nil {
			// Логируем ошибку, но продолжаем, чтобы обновить другие уровни
			continue
		}
	}
	return nil
}

// Delete удаляет ключ из всех уровней кеша.
func (m *MultiLevelCache) Delete(key string) error {
	for _, cache := range m.caches {
		_ = cache.Delete(key) // Игнорируем ошибки
	}
	return nil
}
