package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/github.com/maximiliano745/Geochat-sql/pkg/user"
)

type UserRepository struct {
	Data *Data
}

// Roles
func (ur *UserRepository) GetUserRoles(ctx context.Context, userID uint) ([]user.UserRole, error) {
	roles := []user.UserRole{}

	rows, err := ur.Data.DB.QueryContext(ctx, "SELECT user_id, role_id FROM USER_ROLE WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var role user.UserRole
		if err := rows.Scan(&role.UserID, &role.RolesID); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}

func (ur *UserRepository) SaveUserRole(ctx context.Context, userID uint, roleID uint) error {
	roles, err := ur.GetUserRoles(ctx, userID)
	if err != nil {
		return err
	}

	for _, r := range roles {
		if r.RolesID == roleID {
			fmt.Println("El rol ya existe en este usuario....")
			return nil
		}
	}

	fmt.Print("\nSave User Role...")
	query := `
		INSERT INTO USER_ROLES(user_id, role_id) VALUES($1, $2);
	`
	_, err = ur.Data.DB.ExecContext(ctx, query, userID, roleID)
	if err != nil {
		return err
	}

	return nil
}
func (ur *UserRepository) RemoveUserRole(ctx context.Context, userID uint, roleID uint) error {
	fmt.Print("\nSave User Role...")
	query := `
		DELETE FROM USER_ROLES WHERE user_id = $1 AND role_id = $2;
	`
	_, err := ur.Data.DB.ExecContext(ctx, query, userID, roleID)
	if err != nil {
		return err
	}

	return nil
}

// Trae Miembros Grupo
func (ur *UserRepository) TraeGruposMiembros(ctx context.Context, id uint) ([]user.PartialUser, error) {
	var users []user.User

	// Consulta para obtener los ID de los miembros del grupo dado
	query := `
        SELECT id_miembro FROM grupo_miembros WHERE id_grupo = $1;
    `
	rows, err := ur.Data.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var idMiembro uint
		if err := rows.Scan(&idMiembro); err != nil {
			return nil, err
		}

		// Consulta para obtener los usuarios por ID
		usuario, err := ur.GetOne(ctx, idMiembro)
		if err != nil {
			return nil, err
		}
		users = append(users, usuario)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Convertir los usuarios a PartialUser
	var partialUsers []user.PartialUser
	for _, u := range users {
		partialUser := user.PartialUser{
			ID:       u.ID,
			Username: u.Username,
		}
		partialUsers = append(partialUsers, partialUser)
	}

	fmt.Println("******  Grupo Miembros: ", partialUsers)

	return partialUsers, nil

}

// Trae Grupos
func (ur *UserRepository) TraeGrupos(ctx context.Context, id uint) ([]user.Grupo, error) {
	fmt.Print("\nTraerGrupos------------------------------------------------->")
	var grupos []user.Grupo

	// Consulta SQL para seleccionar los grupos donde el iddueño coincide con el id recibido
	query := `
        SELECT id, group_name FROM user_groups WHERE iddueño = $1;
    `

	rows, err := ur.Data.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var grupo user.Grupo
		if err := rows.Scan(&grupo.ID, &grupo.Nombre); err != nil {
			return nil, err
		}
		grupos = append(grupos, grupo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return grupos, nil
}

// Guardar Grupos
func (ur *UserRepository) CrGrupo(ctx context.Context, g user.Grupo) (int, error) {

	// Verificar si ya existe un grupo con el mismo nombre
	var nombreExistenteID int
	qVerificarNombre := `
        SELECT id FROM user_groups
        WHERE iddueño = $1 AND group_name = $2
    `
	err := ur.Data.DB.QueryRowContext(ctx, qVerificarNombre, g.IDueño, g.Nombre).Scan(&nombreExistenteID)
	if err == nil {
		return nombreExistenteID, errors.New("nombre de grupo existente...... quiere editarlo...??")
	} else if err != sql.ErrNoRows {
		return 0, err
	}
	// Obtener la lista de IDs de miembros para el grupo a crear
	var miembrosIDs []int
	for _, contacto := range g.Contactos {
		miembrosIDs = append(miembrosIDs, int(contacto.ID))
	}

	// Verificar si ya existe un grupo con los mismos miembros
	var grupoExistenteID int
	qVerificar := `
        SELECT id FROM user_groups
        WHERE iddueño = $1
    `
	rows, err := ur.Data.DB.QueryContext(ctx, qVerificar, g.IDueño)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var grupoID int
		err := rows.Scan(&grupoID)
		if err != nil {
			return 0, err
		}

		// Verificar si el grupo actual tiene los mismos miembros
		if gruposIguales(ctx, ur.Data.DB, grupoID, miembrosIDs) {
			grupoExistenteID = grupoID
			break
		}
	}

	if grupoExistenteID != 0 {
		// Si encontramos un grupo existente con los mismos miembros, devolvemos su ID
		return grupoExistenteID, errors.New("ya existe un grupo con los mismos miembros...... quiere editarlo...??")
	}

	// Si no se encontró un grupo existente, procedemos a crear uno nuevo
	q := `
        INSERT INTO user_groups (iddueño, group_name, created_at, updated_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id;
    `
	var grupoID int
	err = ur.Data.DB.QueryRowContext(
		ctx, q, g.IDueño, g.Nombre, time.Now(), time.Now(),
	).Scan(&grupoID)

	if err != nil {
		return 0, err
	}

	fmt.Println("Miembros a guardar: -------------------> ", g.Contactos)
	for _, contacto := range g.Contactos {
		// Suponiendo que el ID del miembro está en contacto.ID
		idMiembro := contacto.ID
		q := `
            INSERT INTO grupo_miembros (id_grupo, id_miembro)
            VALUES ($1, $2)
        `
		_, err := ur.Data.DB.ExecContext(ctx, q, grupoID, idMiembro)
		if err != nil {
			fmt.Println("Error al guardar Miembros del Grupo......", err)
		} else {
			fmt.Println("----------- GUARDADO de Contactos al GRUPO EXITOSO......")
		}
	}

	return grupoID, nil
}

func gruposIguales(ctx context.Context, db *sql.DB, grupoID int, miembrosIDs []int) bool {
	// Obtener la lista de IDs de miembros para el grupo existente
	var miembrosGrupoExistente []int
	q := `
        SELECT id_miembro FROM grupo_miembros
        WHERE id_grupo = $1
    `
	rows, err := db.QueryContext(ctx, q, grupoID)
	if err != nil {
		return false
	}
	defer rows.Close()

	for rows.Next() {
		var miembroID int
		err := rows.Scan(&miembroID)
		if err != nil {
			return false
		}

		miembrosGrupoExistente = append(miembrosGrupoExistente, miembroID)
	}

	// Verificar si las dos listas de miembros son iguales
	return sliceIgual(miembrosIDs, miembrosGrupoExistente)
}

func sliceIgual(slice1 []int, slice2 []int) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}

	return true
}

// Recuperar Contactos
func (ur *UserRepository) GetContactos(ctx context.Context, id uint) ([]int, error) {
	var usuarios []int
	query := `SELECT idusuarioacepta FROM pedidoscontactos WHERE  idusuarioofrece = $1 AND estado = true`

	rows, err := data.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var idusuarioacepta int
		if err := rows.Scan(&idusuarioacepta); err != nil {
			return nil, err
		}
		usuarios = append(usuarios, idusuarioacepta)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return usuarios, nil
}

// Consulta de Pedidos de Contactos (false)
func (ur *UserRepository) ConsultaPedidosContacto() {
	for {
		// Consulta la tabla PedidosContactos para obtener registros con estado en false
		rows, err := data.DB.Query("SELECT idusuarioacepta FROM PedidosContactos WHERE estado = false")
		if err != nil {
			fmt.Println("Error al consultar PedidosContactos:", err)
			time.Sleep(3 * time.Second) // Espera 3 segundos antes de la próxima verificación
			continue
		}

		if err != nil {
			fmt.Println("Error al consultar PedidosContactos:", err)
			time.Sleep(3 * time.Second) // Espera 3 segundos antes de la próxima verificación
			continue
		}
		defer rows.Close()
		// Itera a través de los registros y actualiza el estado según la condición
		for rows.Next() {
			var idusuarioacepta uint
			if err := rows.Scan(&idusuarioacepta); err != nil {
				fmt.Println("Error al escanear fila:", err)
				continue
			}

			// Verifica si el idusuarioacepta existe en la tabla users
			var existe bool
			err := data.DB.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)", idusuarioacepta).Scan(&existe)
			if err != nil {
				fmt.Println("Error al verificar si existe el usuario:", err)
				continue
			}
			// Si el usuario existe, actualiza el estado a true
			if existe {
				_, err := data.DB.Exec("UPDATE PedidosContactos SET estado = true WHERE idusuarioacepta = $1", idusuarioacepta)
				if err != nil {
					fmt.Println("Error al actualizar el estado:", err)
				} else {
					fmt.Printf("Actualizado estado a true para idusuarioacepta = %d\n", idusuarioacepta)
				}
			}
			// Espera 3 segundos antes de la próxima verificación
			time.Sleep(3 * time.Second)
		}

	}
}

// Agrega Pedido de Amistad
func (ur *UserRepository) AgregaPedidoAmistad(ctx context.Context, idusuarioofrece, idusuarioacepta uint) error {
	// Consulta para verificar si ya existe un pedido con los mismos valores
	var pedidoExistente bool
	qExistente := `
        SELECT EXISTS (
            SELECT 1
            FROM PedidosContactos
            WHERE idusuarioofrece = $1 AND idusuarioacepta = $2
        );
    `
	err := ur.Data.DB.QueryRowContext(ctx, qExistente, idusuarioofrece, idusuarioacepta).Scan(&pedidoExistente)
	if err != nil {
		return err
	}

	if pedidoExistente {
		// Si ya existe un pedido idéntico, realiza una actualización en lugar de una inserción
		fmt.Println("\n\n\n Sobreescribiendo Pedido Amistad Existente.....")
		qUpdate := `
            UPDATE PedidosContactos
            SET estado = false -- Puedes agregar otros campos para actualizar si es necesario
            WHERE idusuarioofrece = $1 AND idusuarioacepta = $2;
        `
		_, err := ur.Data.DB.ExecContext(ctx, qUpdate, idusuarioofrece, idusuarioacepta)
		if err != nil {
			return err
		}
	} else {
		// Si no existe un pedido idéntico, realiza una inserción normal
		fmt.Println("\n\n\n Guardando Pedido Amistad.....")
		qInsert := `
            INSERT INTO PedidosContactos (idusuarioofrece, idusuarioacepta)
            VALUES ($1, $2);
        `
		_, err := ur.Data.DB.ExecContext(ctx, qInsert, idusuarioofrece, idusuarioacepta)
		if err != nil {
			return err
		}
	}

	return nil
}

// Obtener todos
func (ur *UserRepository) GetAll(ctx context.Context) ([]user.User, error) {
	q := `
    SELECT id, username, email, pasword, created_at, updated_at
        FROM users;
    `

	rows, err := ur.Data.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []user.User
	for rows.Next() {
		var u user.User
		rows.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
		users = append(users, u)
	}

	return users, nil
}

// Obtener uno
func (ur *UserRepository) GetOne(ctx context.Context, id uint) (user.User, error) {
	q := `
    SELECT id, username, email, pasword, created_at, updated_at
        FROM users WHERE id = $1;
    `

	row := ur.Data.DB.QueryRowContext(ctx, q, id)

	var u user.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return user.User{}, err
	}

	return u, nil
}

// Obtener por username
func (ur *UserRepository) GetByUsername(ctx context.Context, username string) (user.User, error) {
	q := `
    SELECT id,  username, email, password, created_at, updated_at
        FROM users WHERE username = $1;
    `

	row := ur.Data.DB.QueryRowContext(ctx, q, username)

	var u user.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Hash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return user.User{}, err
	}

	return u, nil
}

// Obtener X Email  ----------------- MIO -----------------------
func (ur *UserRepository) GetByMail(ctx context.Context, email string) (user.User, error) {
	fmt.Println("Revisanso Email existe..???   GetByMail()")
	q := `
    SELECT  id, username, email, hash, created_at, updated_at, pasword
        FROM users WHERE email = $1;
    `

	row := ur.Data.DB.QueryRowContext(ctx, q, email)

	var u user.User
	err := row.Scan(&u.ID, &u.Username, &u.Email,
		&u.Hash, &u.CreatedAt, &u.UpdatedAt, &u.Password)
	if err != nil {
		return user.User{}, err
	}

	return u, nil
}

// Insertar
func (ur *UserRepository) Create(ctx context.Context, u *user.User) error {
	q := `
    INSERT INTO users (username, pasword, email, created_at, updated_at, hash)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id;
    `

	if err := u.HashPassword(); err != nil {
		return err
	}

	row := ur.Data.DB.QueryRowContext(
		ctx, q, u.Username, u.Password, u.Email, time.Now(), time.Now(), u.Hash,
	)

	err := row.Scan(&u.ID)
	if err != nil {
		return err
	}

	return nil
}

// Actualizar
func (ur *UserRepository) Update(ctx context.Context, id uint, u user.User) error {
	q := `
    UPDATE users set username=$1, email=$2, updated_at=$3
        WHERE id=$3;
    `

	stmt, err := ur.Data.DB.PrepareContext(ctx, q)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx, u.Email, time.Now(), id,
	)
	if err != nil {
		return err
	}

	return nil
}

// Borrar
func (ur *UserRepository) Delete(ctx context.Context, id uint) error {
	q := `DELETE FROM users WHERE id=$1;`

	stmt, err := ur.Data.DB.PrepareContext(ctx, q)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
