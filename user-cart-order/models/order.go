package models

import "time"

type Order struct {
	ID        int
	UserID    int
	ProductID int
	CreatedAt time.Time
	Status    string
}
