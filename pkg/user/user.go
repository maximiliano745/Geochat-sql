package user

import (
	"time"

	//"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type GrupoMiembros struct {
	NombreGrupo    string     `json:"nombregrupo"`
	IdGrupo        uint       `json:"idgrupo"`
	Iddueño        uint       `json:"iddueño"`
	ContactosGrupo []Contacto `json:"contactosgrupo"`
}

type Grupo struct {
	Nombre    string     `json:"nombre"`
	ID        uint       `json:"id"`
	IDueño    uint       `json:"iddueño"`
	Contactos []Contacto `json:"contactos"`
}

type Contacto struct {
	ID     uint   `json:"id"`
	Nombre string `json:"nombre"`
}

type PartialUser struct {
	ID       uint   `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
}

type User struct {
	ID        uint      `json:"id,omitempty"`
	Username  string    `json:"username,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"Password,omitempty"`
	Hash      string    `json:"hash,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type UserRole struct {
	UserID  uint `json:"user_id,omitempty"`
	RolesID uint `json:"roles_id,omitempty"`
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
