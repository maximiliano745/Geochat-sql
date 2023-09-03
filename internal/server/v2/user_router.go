package v2

import (
	//"context"

	"encoding/json"
	"fmt"

	//"log"
	"net/http"

	//"time"

	//"strconv"

	"github.com/github.com/maximiliano745/Geochat-sql/pkg/response"
	"github.com/github.com/maximiliano745/Geochat-sql/pkg/user"
	"github.com/go-chi/chi"
	"golang.org/x/crypto/bcrypt"
)

type UserRouter struct {
	Repository user.Repository
}

var SECRET_KEY = []byte("gosecretkey")

/* func (ur *UserRouter) HandleFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Bienvenidos a Geochat Backend con SQL")
} */

func getHash(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (ur *UserRouter) UserLogin(w http.ResponseWriter, r *http.Request) {
	var u, uu user.User
	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	uu = u
	defer r.Body.Close()

	fmt.Println("**** ACA ESTAMOS EN EL LOGIN   ******")

	ctx := r.Context()
	uu, err = ur.Repository.GetByMail(ctx, uu.Email)
	if err == nil {
		fmt.Println(" OK!!, Email EXISTENTE....!!!")
		w.Write([]byte(`{"OK!!, Email EXISTENTE....!!!"}`))
		//w.WriteHeader(http.StatusOK)
		fmt.Println("Nombre:", uu.Username, "     ID:", uu.ID, "      Password:", uu.Password, "    Email:", uu.Email, "   First_Name:", uu.FirstName, "   Last_Name:", uu.LastName)
		fmt.Println("Revisando Contraseña.....", u.Password, "     ", uu.Hash)

		if uu.PasswordMatch(u.Password) {
			fmt.Println("Contraseña correcta")
		} else {
			fmt.Println("ERROR Contraseña INCORECTA...!!!")
			w.Write([]byte(`{"ERROR Contraseña INCORECTA...!!!"}`))
			return
		}
	} else {
		fmt.Println("ERROR Email Inexistente....!!!")
		w.Write([]byte(`{"ERROR Email Inexistente....!!!"}`))
		//w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (ur *UserRouter) UserSignup(w http.ResponseWriter, r *http.Request) {
	var u, uu user.User
	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	uu = u
	defer r.Body.Close()

	ctx := r.Context()
	uu, err = ur.Repository.GetByMail(ctx, uu.Email)
	if err == nil {
		fmt.Println("Error Email EXISTENTE....!!!")
		w.Write([]byte(`{"Error Email EXISTENTE....!!!"}`))
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		uu = u
		u.Hash, err = getHash([]byte(u.Password))
		if err != nil {
			fmt.Println("Error Al Generar Hash")
			//response.HTTPError(w, r, http.StatusUnauthorized, err.Error())
			w.Write([]byte(`{"ERROR Al Generar Hash"}`))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = ur.Repository.Create(ctx, &u)
		if err != nil {
			fmt.Println("\nERROR al Guardar: ", err)
			w.Write([]byte(`{"ERROR al Guardar"}`))
			w.WriteHeader(http.StatusBadRequest)
			//response.HTTPError(w, r, http.StatusBadRequest, err.Error())
			return
		} else {
			fmt.Println("USUARIO CREADO CON  EXITO...!!!!: ")
			fmt.Println("Nombre:", u.Username, "     ID:", u.ID, "      Password:", u.Password, "    Email:", u.Email, "   First_Name:", u.FirstName, "   Last_Name:", u.LastName)
			w.Write([]byte(`{"USUARIO CREADO CON  EXITO...!!!!"}`))
			w.WriteHeader(http.StatusOK)
		}
	}
}

// ****************     Definiendo rutas    ************************
func (ur *UserRouter) Routes() http.Handler {
	r := chi.NewRouter()

	// Configurar el middleware CORS para permitir todas las solicitudes desde cualquier origen

	//r.Get("/", ur.HandleFunc)
	r.Post("/login", ur.UserLogin)
	r.Post("/", ur.UserSignup) // http://localhost:9000/api/v2/users/
	//r.Get("/{id}", ur.GetOneHandler)
	//r.Put("/{id}", ur.UpdateHandler)
	//r.Delete("/{id}", ur.DeleteHandler)

	return r
}
