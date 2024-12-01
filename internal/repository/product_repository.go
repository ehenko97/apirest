package repository

import (
	"Projectapirest/internal/entity"
	"context"
	"database/sql"
	"time"
)

// ProductRepositoryInterface описывает методы работы с продуктами.
type ProductRepositoryInterface interface {
	Create(ctx context.Context, product entity.Product) (entity.Product, error)
	FindByID(ctx context.Context, id int) (entity.Product, error)
	Update(ctx context.Context, product entity.Product) error
	Delete(ctx context.Context, id int) error
	FindAll(ctx context.Context) ([]entity.Product, error)
}

// ProductRepository содержит ссылку на базу данных и реализует интерфейс ProductRepositoryInterface.
type ProductRepository struct {
	db *sql.DB
}

// NewProductRepository создает новый репозиторий продуктов.
func NewProductRepository(db *sql.DB) ProductRepositoryInterface {
	return &ProductRepository{db: db} // Возвращаем интерфейс
}

// Create добавляет новый продукт в базу данных.
func (r *ProductRepository) Create(ctx context.Context, product entity.Product) (entity.Product, error) {
	query := `
        INSERT INTO product (name, description, price, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRowContext(
		ctx,
		query,
		product.Name,
		product.Description,
		product.Price,
		time.Now(),
		time.Now(),
	).Scan(&product.ID)
	if err != nil {
		return product, err
	}
	return product, nil
}

// FindByID находит продукт по ID.
func (r *ProductRepository) FindByID(ctx context.Context, id int) (entity.Product, error) {
	query := `SELECT id, name, description, price, created_at, updated_at FROM product WHERE id = $1`
	var product entity.Product
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return product, err
	}
	return product, nil
}

// Update обновляет информацию о продукте.
func (r *ProductRepository) Update(ctx context.Context, product entity.Product) error {
	query := `
        UPDATE product
        SET name = $1, description = $2, price = $3, updated_at = $4
        WHERE id = $5`
	_, err := r.db.ExecContext(
		ctx,
		query,
		product.Name,
		product.Description,
		product.Price,
		time.Now(),
		product.ID,
	)
	return err
}

// Delete удаляет продукт из базы данных.
func (r *ProductRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM product WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// FindAll возвращает список всех продуктов.
func (r *ProductRepository) FindAll(ctx context.Context) ([]entity.Product, error) {
	query := `SELECT id, name, description, price, created_at, updated_at FROM product`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []entity.Product{}
	for rows.Next() {
		var product entity.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}
