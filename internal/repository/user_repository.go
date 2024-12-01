package repository

import (
	"Projectapirest/internal/entity"
	"context"
	"database/sql"
	"time"
)

// UserRepositoryInterface описывает методы работы с пользователями.
type UserRepositoryInterface interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	FindByID(ctx context.Context, id int) (entity.User, error)
	Update(ctx context.Context, user entity.User) error
	Delete(ctx context.Context, id int) error
	FindAll(ctx context.Context) ([]entity.User, error)
}

// UserRepository содержит ссылку на базу данных и реализует интерфейс UserRepositoryInterface.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository создает новый репозиторий пользователей.
func NewUserRepository(db *sql.DB) UserRepositoryInterface {
	return &UserRepository{db: db} // Возвращаем интерфейс
}

// Create добавляет нового пользователя в базу данных.
func (r *UserRepository) Create(ctx context.Context, user entity.User) (entity.User, error) {
	query := `
        INSERT INTO users (name, email, password, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Name,
		user.Email,
		time.Now(),
		time.Now(),
	).Scan(&user.ID)
	if err != nil {
		return user, err
	}
	return user, nil
}

// FindByID находит пользователя по ID.
func (r *UserRepository) FindByID(ctx context.Context, id int) (entity.User, error) {
	query := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE id = $1`
	var user entity.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return user, err
	}
	return user, nil
}

// Update обновляет информацию о пользователе.
func (r *UserRepository) Update(ctx context.Context, user entity.User) error {
	query := `
        UPDATE users
        SET name = $1, email = $2, password = $3, updated_at = $4
        WHERE id = $5`
	_, err := r.db.ExecContext(
		ctx,
		query,
		user.Name,
		user.Email,
		time.Now(),
		user.ID,
	)
	return err
}

// Delete удаляет пользователя из базы данных.
func (r *UserRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// FindAll возвращает список всех пользователей.
func (r *UserRepository) FindAll(ctx context.Context) ([]entity.User, error) {
	query := `SELECT id, name, email, password, created_at, updated_at FROM users`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []entity.User{}
	for rows.Next() {
		var user entity.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
