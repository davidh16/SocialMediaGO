package Models

import (
	"time"
)

type Post struct {
	ID          int   `gorm:"primarykey"`
	UserId      int   `json:"-"`
	User        *User `gorm:"constraint:OnDelete:CASCADE"`
	Description string
	Image       string
	Deleted     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `gorm:"default:"`
}
