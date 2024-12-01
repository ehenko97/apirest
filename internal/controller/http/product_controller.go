package http

import (
	"Projectapirest/internal/entity"
	service "Projectapirest/internal/services"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

// ProductController структура контроллера для обработки запросов
type ProductController struct {
	productService service.ProductService
}

// NewProductController создает новый контроллер продуктов
func NewProductController(productService service.ProductService) *ProductController {
	return &ProductController{
		productService: productService,
	}
}

// CreateProduct создает новый продукт
func (pc *ProductController) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product entity.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdProduct, err := pc.productService.Create(r.Context(), product)
	if err != nil {
		http.Error(w, "Failed to create product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdProduct)
}

// GetAllProducts возвращает все продукты
func (pc *ProductController) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := pc.productService.FindAll(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch products: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(products)
}

// GetProduct возвращает продукт по ID
func (pc *ProductController) GetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := extractProductIDFromURL(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid product ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	product, err := pc.productService.FindByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Product not found: "+err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(product)
}

// UpdateProduct обновляет продукт по ID
func (pc *ProductController) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := extractProductIDFromURL(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid product ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	var updatedProduct entity.Product
	if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	updatedProduct.ID = id

	product, err := pc.productService.Update(r.Context(), updatedProduct)
	if err != nil {
		http.Error(w, "Failed to update product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(product)
}

// DeleteProduct удаляет продукт по ID
func (pc *ProductController) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := extractProductIDFromURL(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid product ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = pc.productService.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to delete product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// extractIDFromURL извлекает ID из URL
func extractProductIDFromURL(path string) (int, error) {
	parts := strings.Split(strings.TrimSuffix(path, "/"), "/")
	if len(parts) < 5 { // Проверка, что путь содержит хотя бы 5 частей
		return 0, errors.New("ID is missing in the URL")
	}
	return strconv.Atoi(parts[4]) // ID теперь в 5-й части пути
}
