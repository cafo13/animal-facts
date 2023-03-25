package database

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Database
	gorm.Model
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null;" json:"password"`
}

func (u *User) CheckLogin() error {
	password := u.Password
	err := u.Database.db.Model(User{}).Where("username = ?", u.Username).Take(&u).Error
	if err != nil {
		return errors.New("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return errors.New("invalid username or password")
	}

	return nil
}
