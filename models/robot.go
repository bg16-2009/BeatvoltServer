package models

import (
	"golang.org/x/crypto/bcrypt"
)

type Robot struct {
	ID                 uint   `gorm:"primaryKey"`
	Username           string `gorm:"uniqueIndex;not null"`
	Password           string `gorm:"not null"`
	CollectedBatteries int    `gorm:"default:0"`
}

func (r *Robot) SetPassword(rawPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	r.Password = string(hashedPassword)
	return nil
}

func (r *Robot) CheckPassword(rawPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(r.Password), []byte(rawPassword)) == nil
}
