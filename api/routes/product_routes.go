package routes

import (
	http2 "Projectapirest/internal/controller/http"
	"net/http"
)

// SetupProductRoutes настраивает маршруты для работы с продуктами
func SetupProductRoutes(productController *http2.ProductController) *http.ServeMux {
	mux := http.NewServeMux()

	// Маршруты для работы с продуктами
	mux.HandleFunc("/api/v1/products", productController.GetAllProducts) // GET для всех продуктов
	mux.HandleFunc("/api/v1/products", productController.CreateProduct)  // POST для создания продукта
	mux.HandleFunc("/api/v1/products/", productController.GetProduct)    // GET для продукта по ID
	mux.HandleFunc("/api/v1/products/", productController.UpdateProduct) // PUT для обновления продукта
	mux.HandleFunc("/api/v1/products/", productController.DeleteProduct) // DELETE для удаления продукта

	return mux
}
