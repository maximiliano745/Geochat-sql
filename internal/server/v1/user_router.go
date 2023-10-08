package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/github.com/maximiliano745/Geochat-sql/pkg/response"
	"github.com/github.com/maximiliano745/Geochat-sql/pkg/user"
	"github.com/go-chi/chi"
)

type UserRouter struct {
	Repository user.Repository
}

// Crear Usuario
func (ur *UserRouter) CreateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Aca estamos en CREAR Usuario\n")
	var u user.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	ctx := r.Context()
	fmt.Print("\n---------------- > Pasword: " + u.Password)
	err = ur.Repository.Create(ctx, &u)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	u.Password = ""
	w.Header().Add("Location", fmt.Sprintf("%s%d", r.URL.String(), u.ID))
	response.JSON(w, r, http.StatusCreated, response.Map{"user": u})
}

// Obtener usuarios
func (ur *UserRouter) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n\n\nAca estamos en el Get de Maxi.....")
	ctx := r.Context()

	users, err := ur.Repository.GetAll(ctx)
	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, response.Map{"users": users})
}

// Obtener un usuario
func (ur *UserRouter) GetOneHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	u, err := ur.Repository.GetOne(ctx, uint(id))
	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, response.Map{"user": u})
}

// Actualizar usuario
func (ur *UserRouter) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	var u user.User
	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	ctx := r.Context()
	err = ur.Repository.Update(ctx, uint(id), u)
	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, nil)
}

// Borrar un usuarios
func (ur *UserRouter) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	err = ur.Repository.Delete(ctx, uint(id))
	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, response.Map{})
}

// ****************     Definiendo rutas    ************************
func (ur *UserRouter) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", ur.GetAllHandler) //  http://localhost:9000/api/v1/Maxi/  GET
	r.Post("/", ur.CreateHandler)
	r.Get("/{id}", ur.GetOneHandler)
	r.Put("/{id}", ur.UpdateHandler)
	r.Delete("/{id}", ur.DeleteHandler)

	return r
}
