package routes

import (
	http2 "Projectapirest/internal/controller/http"
	"net/http"
)

// SetupUserRoutes настраивает маршруты для работы с пользователями
func SetupUserRoutes(userController *http2.UserController) *http.ServeMux {
	mux := http.NewServeMux()

	// Маршруты для работы с пользователями
	mux.HandleFunc("/api/v1/users", userController.GetAllUsers) // GET для всех пользователей
	mux.HandleFunc("/api/v1/users", userController.CreateUser)  // POST для создания пользователя
	mux.HandleFunc("/api/v1/users/", userController.GetUser)    // GET для пользователя по ID
	mux.HandleFunc("/api/v1/users/", userController.UpdateUser) // PUT для обновления пользователя
	mux.HandleFunc("/api/v1/users/", userController.DeleteUser) // DELETE для удаления пользователя

	return mux
}
