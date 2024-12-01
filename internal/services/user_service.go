package service

import (
	"Projectapirest/internal/cache"
	"Projectapirest/internal/entity"
	"Projectapirest/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// UserService описывает методы управления пользователями.
type UserService interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	FindByID(ctx context.Context, id int) (entity.User, error)
	Update(ctx context.Context, user entity.User) (entity.User, error)
	Delete(ctx context.Context, id int) error
	FindAll(ctx context.Context) ([]entity.User, error)
}

// DecodeUserRequestBody десериализует тело запроса в структуру.
func DecodeUserRequestBody(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(v)
}

// EncodeUserResponse сериализует структуру в JSON и отправляет ее в ответ.
func EncodeUserResponse(w http.ResponseWriter, v interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(v)
}

type userService struct {
	repo  repository.UserRepository
	cache cache.Cache
}

func NewUserService(repo repository.UserRepository, cache cache.Cache) UserService {
	return &userService{
		repo:  repo,
		cache: cache,
	}
}

// Create добавляет нового пользователя.
func (s *userService) Create(ctx context.Context, user entity.User) (entity.User, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	createdUser, err := s.repo.Create(ctx, user)
	if err != nil {
		return entity.User{}, err
	}

	// Инвалидация кеша списка пользователей
	_ = s.cache.Delete("users:all")

	return createdUser, nil
}

// FindByID находит пользователя по ID.
func (s *userService) FindByID(ctx context.Context, id int) (entity.User, error) {
	cacheKey := fmt.Sprintf("user:%d", id)

	// Попытка извлечь из кеша
	if cached, err := s.cache.Get(cacheKey); err == nil && cached != "" {
		var user entity.User
		if err := json.Unmarshal([]byte(cached), &user); err == nil {
			return user, nil
		}
	}

	// Если не найдено в кеше, обращаемся к репозиторию
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return entity.User{}, errors.New("пользователь не найден")
	}

	// Сохраняем результат в кеш
	userJSON, _ := json.Marshal(user)
	if err := s.cache.Set(cacheKey, string(userJSON), 60); err != nil {
		fmt.Printf("Failed to set cache for user %d: %v\n", id, err)
	}

	return user, nil
}

// Update обновляет информацию о пользователе.
func (s *userService) Update(ctx context.Context, user entity.User) (entity.User, error) {
	user.UpdatedAt = time.Now()

	err := s.repo.Update(ctx, user)
	if err != nil {
		return entity.User{}, err
	}

	// Инвалидация кеша пользователя и списка пользователей
	cacheKey := fmt.Sprintf("user:%d", user.ID)
	_ = s.cache.Delete(cacheKey)
	_ = s.cache.Delete("users:all")

	return user, nil
}

// Delete удаляет пользователя.
func (s *userService) Delete(ctx context.Context, id int) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Удаляем из кеша пользователя и списка пользователей
	cacheKey := fmt.Sprintf("user:%d", id)
	_ = s.cache.Delete(cacheKey)
	_ = s.cache.Delete("users:all")

	return nil
}

// FindAll возвращает всех пользователей.
func (s *userService) FindAll(ctx context.Context) ([]entity.User, error) {
	cacheKey := "users:all"

	// Попытка извлечь данные из кеша
	if cached, err := s.cache.Get(cacheKey); err == nil && cached != "" {
		var users []entity.User
		if err := json.Unmarshal([]byte(cached), &users); err == nil {
			return users, nil
		}
	}

	// Если кеш пуст, извлекаем данные из репозитория
	users, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// Сохраняем результат в кеш
	usersJSON, _ := json.Marshal(users)
	if err := s.cache.Set(cacheKey, string(usersJSON), 300); err != nil {
		fmt.Printf("Failed to set cache for FindAll: %v\n", err)
	}

	return users, nil
}
