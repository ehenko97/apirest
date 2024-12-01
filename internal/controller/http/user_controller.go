package http

import (
	"Projectapirest/internal/entity"
	service "Projectapirest/internal/services"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

// UserController структура контроллера для обработки запросов.
type UserController struct {
	userService service.UserService
}

// NewUserController создает новый контроллер пользователей.
func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// CreateUser создает нового пользователя.
func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	// Десериализация запроса через функцию из сервиса
	if err := service.DecodeUserRequestBody(r, &user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdUser, err := uc.userService.Create(r.Context(), user)
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Сериализация ответа через функцию из сервиса
	service.EncodeUserResponse(w, createdUser, http.StatusCreated)
}

// GetAllUsers возвращает всех пользователей.
func (uc *UserController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := uc.userService.FindAll(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Сериализация ответа через функцию из сервиса
	service.EncodeUserResponse(w, users, http.StatusOK)
}

// GetUser возвращает пользователя по ID.
func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := extractUserIDFromURL(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid user ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := uc.userService.FindByID(r.Context(), id)
	if err != nil {
		http.Error(w, "User not found: "+err.Error(), http.StatusNotFound)
		return
	}

	// Сериализация ответа через функцию из сервиса
	service.EncodeUserResponse(w, user, http.StatusOK)
}

// UpdateUser обновляет пользователя по ID.
func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := extractUserIDFromURL(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid user ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	var updatedUser entity.User
	// Десериализация запроса через функцию из сервиса
	if err := service.DecodeUserRequestBody(r, &updatedUser); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	updatedUser.ID = id

	user, err := uc.userService.Update(r.Context(), updatedUser)
	if err != nil {
		http.Error(w, "Failed to update user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Сериализация ответа через функцию из сервиса
	service.EncodeUserResponse(w, user, http.StatusOK)
}

// DeleteUser удаляет пользователя по ID.
func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := extractUserIDFromURL(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid user ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = uc.userService.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to delete user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// extractUserIDFromURL извлекает ID пользователя из URL.
func extractUserIDFromURL(path string) (int, error) {
	parts := strings.Split(strings.TrimSuffix(path, "/"), "/")
	if len(parts) < 5 { // Проверка, что путь содержит хотя бы 5 частей
		return 0, errors.New("ID is missing in the URL")
	}
	return strconv.Atoi(parts[4]) // ID теперь в 5-й части пути
}
