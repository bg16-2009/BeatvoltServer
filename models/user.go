package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                uint   `gorm:"primaryKey"`
	Username          string `gorm:"uniqueIndex;not null"`
	Email             string `gorm:"uniqueIndex;not null"`
	Password          string `gorm:"not null"`
	IsAdmin           bool   `gorm:"default:false"`
	RecycledBatteries int    `gorm:"default:0"`
}

func (u *User) SetPassword(rawPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(rawPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(rawPassword)) == nil
}
