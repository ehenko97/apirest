package service

import (
	"Projectapirest/internal/cache"
	"Projectapirest/internal/entity"
	"Projectapirest/internal/repository"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ProductService описывает методы управления продуктами.
type ProductService interface {
	Create(ctx context.Context, product entity.Product) (entity.Product, error)
	FindByID(ctx context.Context, id int) (entity.Product, error)
	Update(ctx context.Context, product entity.Product) (entity.Product, error)
	Delete(ctx context.Context, id int) error
	FindAll(ctx context.Context) ([]entity.Product, error)
}

// DecodeRequestBody десериализует тело запроса в структуру.
func DecodeRequestBody(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(v)
}

// EncodeResponse сериализует структуру в JSON и отправляет ее в ответ.
func EncodeResponse(w http.ResponseWriter, v interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(v)
}

// productService реализует интерфейс ProductService.
type productService struct {
	repo  repository.ProductRepositoryInterface
	cache cache.Cache
}

// NewProductService создает новый экземпляр productService.
func NewProductService(repo repository.ProductRepositoryInterface, inMemoryCache, redisCache cache.Cache) ProductService {
	// Создаем многослойный кеш
	multiCache := cache.NewMultiLevelCache(inMemoryCache, redisCache)

	return &productService{
		repo:  repo,
		cache: multiCache,
	}
}

// Create добавляет новый продукт.
func (s *productService) Create(ctx context.Context, product entity.Product) (entity.Product, error) {
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	createdProduct, err := s.repo.Create(ctx, product)
	if err != nil {
		return entity.Product{}, err
	}

	// Инвалидация кеша списка продуктов
	_ = s.cache.Delete("products:all")

	return createdProduct, nil
}

// FindByID находит продукт по ID.
func (s *productService) FindByID(ctx context.Context, id int) (entity.Product, error) {
	cacheKey := fmt.Sprintf("product:%d", id)

	// Попытка извлечь из кеша
	if cached, err := s.cache.Get(cacheKey); err == nil && cached != "" {
		var product entity.Product
		if err := json.Unmarshal([]byte(cached), &product); err == nil {
			return product, nil
		}
	}

	// Если кеш пуст, получаем из базы
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return entity.Product{}, err
	}

	// Сохраняем результат в кеш
	productJSON, _ := json.Marshal(product)
	if err := s.cache.Set(cacheKey, string(productJSON), 300); err != nil {
		fmt.Printf("Failed to set cache for FindByID: %v\n", err)
	}

	return product, nil
}

// Update обновляет информацию о продукте.
func (s *productService) Update(ctx context.Context, product entity.Product) (entity.Product, error) {
	product.UpdatedAt = time.Now()

	err := s.repo.Update(ctx, product)
	if err != nil {
		return entity.Product{}, err
	}

	// Инвалидация кеша продукта и списка продуктов
	_ = s.cache.Delete(fmt.Sprintf("product:%d", product.ID))
	_ = s.cache.Delete("products:all")

	return product, nil
}

// Delete удаляет продукт.
func (s *productService) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Инвалидация кеша продукта и списка продуктов
	_ = s.cache.Delete(fmt.Sprintf("product:%d", id))
	_ = s.cache.Delete("products:all")

	return nil
}

// FindAll возвращает всех продуктов.
func (s *productService) FindAll(ctx context.Context) ([]entity.Product, error) {
	cacheKey := "products:all"

	// Попытка извлечь данные из кеша
	if cached, err := s.cache.Get(cacheKey); err == nil && cached != "" {
		var products []entity.Product
		if err := json.Unmarshal([]byte(cached), &products); err == nil {
			return products, nil
		}
	}

	// Если кеш пуст, извлекаем данные из репозитория
	products, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// Сохраняем результат в кеш
	productsJSON, _ := json.Marshal(products)
	if err := s.cache.Set(cacheKey, string(productsJSON), 300); err != nil {
		fmt.Printf("Failed to set cache for FindAll: %v\n", err)
	}

	return products, nil
}
