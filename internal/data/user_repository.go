package data

import (
	"context"
	"fmt"
	"time"

	"github.com/github.com/maximiliano745/Geochat-sql/pkg/user"
)

type UserRepository struct {
	Data *Data
}

// Consulta de Pedidos de Contactos (false)
func (ur *UserRepository) ConsultaPedidosContacto() {
	for {
		//fmt.Println("\n\n\n ")
		// Coloca aquí el código que deseas que se ejecute cada 3 segundos.
		//fmt.Println("Verificando y actualizando PedidosContactos...")

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
    SELECT id, username, email, password, created_at, updated_at
        FROM users WHERE id = $1;
    `

	row := ur.Data.DB.QueryRowContext(ctx, q, id)

	var u user.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt, &u.UpdatedAt)
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
