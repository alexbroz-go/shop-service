package models

import "time"

type Cart struct {
	ID        int
	UserID    int
	ProductID int
	CreatedAt time.Time
}
