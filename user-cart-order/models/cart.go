package models

import "time"

type Cart struct {
	ID        int
	UserID    int
	CreatedAt time.Time
}
