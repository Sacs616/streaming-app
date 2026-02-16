package domain

import (
    "time"
    "golang.org/x/crypto/bcrypt"
)

type User struct {
    ID        uint64    `gorm:"primaryKey"`
    Email     string    `gorm:"uniqueIndex;not null"`
    Username  string    `gorm:"uniqueIndex;not null"`
    Password  string    `gorm:"not null"`
    Role      string    `gorm:"default:'user'"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (u *User) HashPassword(password string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    u.Password = string(hashedPassword)
    return nil
}

func (u *User) CheckPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
    return err == nil
}