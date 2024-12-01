package entity

import "time"

// Струтура определяет какие данные будут храниться в базе данных.

type User struct {
	ID        int
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
