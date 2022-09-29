package Models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID          int `gorm:"primarykey"`
	Name        string
	Surname     string
	Email       string `gorm:"unique" json:"-"`
	Password    string `json:"-"`
	Verified    bool   `gorm:"default:false"`
	Deactivated bool   `gorm:"default:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `gorm:"default:"`
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return
}
