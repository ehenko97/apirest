package entity

import "time"

// Струтура определяет какие данные будут храниться в базе данных.

type Product struct {
	ID          int
	Name        string
	Description string
	Price       float64
	UserID      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
