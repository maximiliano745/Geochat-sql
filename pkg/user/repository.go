package user

import (
	"context"
)

// Repository handle the CRUD operations with Users.
type Repository interface {
	GetAll(ctx context.Context) ([]User, error) // Maxi
	GetOne(ctx context.Context, id uint) (User, error)

	GetByUsername(ctx context.Context, username string) (User, error)

	GetByMail(ctx context.Context, email string) (User, error) // ------------- MIO --------------
	AgregaPedidoAmistad(ctx context.Context, emailOfrece, emailAcepta uint) error
	ConsultaPedidosContacto()
	GetContactos(ctx context.Context, id uint) ([]int, error)
	CrGrupo(ctx context.Context, g Grupo) (int, error)
	TraeGrupos(ctx context.Context, id uint) ([]Grupo, error)
	TraeGruposMiembros(ctx context.Context, id uint) ([]PartialUser, error)

	// Roles
	GetUserRoles(ctx context.Context, userID uint) ([]UserRole, error)
	SaveUserRole(ctx context.Context, userID, roleID uint) error
	RemoveUserRole(ctx context.Context, userID uint, roleID uint) error

	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, id uint, user User) error
	Delete(ctx context.Context, id uint) error
}
