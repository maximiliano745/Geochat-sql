package v3

import (
	//"context"

	"encoding/json"
	"fmt"
	"time"

	//"log"
	"net/http"

	//"time"

	//"strconv"

	"github.com/github.com/maximiliano745/Geochat-sql/pkg/response"
	"github.com/github.com/maximiliano745/Geochat-sql/pkg/user"
	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt"
)

type UserRouter struct {
	Repository user.Repository
}

var SECRET_KEY = []byte("gosecretkey")

func (ur *UserRouter) LoginMovil(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	fmt.Println("**** ACA ESTAMOS EN EL LOGIN NATIVO  ******")

	var status bool
	var msg string

	fmt.Println("\n -------------- Aca estamos en el Login Nativo. ---------------- ")

	if r.Method != http.MethodPost {
		msg = "Error metodo request"
		http.Error(w, "Error metodo request", http.StatusMethodNotAllowed)
		return
	}

	var u, uu user.User
	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println("Email:", uu.Email)
	fmt.Println("Password:", uu.Password)

	// Incluir el correo electrónico en las reclamaciones
	claims := jwt.MapClaims{
		"email": uu.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Caducidad del token
	}

	// Generar el token JWT con las reclamaciones
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	mySecret := "secret-secret"
	signedToken, err := token.SignedString([]byte(mySecret))
	if err != nil {
		msg = "Error al generar el token"
		http.Error(w, "Error al generar el token", http.StatusInternalServerError)
		return
	} else {
		fmt.Println("signedToken :" + signedToken)
		//fmt.Println("Token: ", token)
	}

	responseData := map[string]interface{}{

		"status": status,
		"msg":    msg,
		"token":  signedToken, // Envía el token firmado en la respuesta
	}

	// Convertir el mapa a formato JSON
	jsonResponse, err := json.Marshal(responseData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		msg = "Error al generar respuesta JSON"
		http.Error(w, "Error al generar respuesta JSON", http.StatusInternalServerError)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("aca claims: ", claims["email"], claims["exp"])
	}

	// Establecer la cabecera Content-Type y enviar la respuesta JSON al cliente
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

/* func (ur *UserRouter) UserLogin(w http.ResponseWriter, r *http.Request) {
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

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("----------------  WebSocket Endpoint Hit -------------------")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
} */

// ****************     Definiendo rutas    ************************
func (ur *UserRouter) Routes() http.Handler {
	r := chi.NewRouter()

	//r.Post("/login", ur.UserLogin)

	r.Post("/movil", ur.LoginMovil) //    fetch('https://geochat-efn9.onrender.com/api/v3/users/movil')
	return r
}
