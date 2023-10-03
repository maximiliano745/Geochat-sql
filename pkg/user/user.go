package user

import (
	"time"

	//"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// User of the system.
type User struct {
	ID        uint      `json:"id,omitempty"`
	Username  string    `json:"username,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"Password,omitempty"`
	Hash      string    `json:"hash,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (u *User) HashPassword() error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Hash = string(passwordHash)

	return nil
}

func (u User) PasswordMatch(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Hash), []byte(password))

	return err == nil
}
