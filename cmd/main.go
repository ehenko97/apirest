package main

import (
	http2 "Projectapirest/internal/controller/http"
	"Projectapirest/internal/repository"
	"database/sql"
	"log"
	"net/http"

	"Projectapirest/api/routes"
	"Projectapirest/internal/services"
	_ "github.com/lib/pq" // Импорт драйвера для PostgreSQL
)

func main() {

	// Настройка подключения к базе данных
	db, err := sql.Open("postgres", "postgresql://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Проверка соединения с базой
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	// Инициализация репозитория и сервиса
	productRepo := repository.NewProductRepository(db)              // Создаем репозиторий (тип *ProductRepository)
	productService := service.NewProductService(productRepo)        // Передаем *ProductRepository в сервис
	productController := http2.NewProductController(productService) // Создаем контроллер с сервисом

	// Настройка маршрутов
	mux := routes.SetupProductRoutes(productController)

	// Запуск сервера
	port := ":8080"
	log.Printf("Server is running on port %s", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
