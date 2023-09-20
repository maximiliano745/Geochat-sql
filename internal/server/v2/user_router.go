package v2

import (
	//"context"

	"encoding/json"
	"fmt"
	"log"

	//"log"
	"net/http"
	"net/smtp"

	//"time"

	//"strconv"

	"github.com/github.com/maximiliano745/Geochat-sql/pkg/response"
	"github.com/github.com/maximiliano745/Geochat-sql/pkg/user"
	"github.com/github.com/maximiliano745/Geochat-sql/pkg/websocket"
	"github.com/go-chi/chi"
	"golang.org/x/crypto/bcrypt"
)

type UserRouter struct {
	Repository user.Repository
}

var SECRET_KEY = []byte("gosecretkey")

func getHash(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (ur *UserRouter) UserMail(w http.ResponseWriter, r *http.Request) {
	var u user.User
	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	fmt.Println("**** ACA ESTAMOS EN EL EMAIL   ******")

	auth := smtp.PlainAuth("", "maxiargento745@gmail.com", "rwkycxemzftxidxi", "smtp.gmail.com")

	to := []string{u.Email}
	msg := []byte("To: " + u.Email + "\r\n" +
		"Subject: Geochat..!!!\r\n" +
		"\r\n" +
		"Esto es la Invitacion de Contacto de GEOCHAT  ---------------->   " + "https://maxi-geochat.onrender.com/")
	err = smtp.SendMail("smtp.gmail.com:587", auth, "maxiargento745@gmail.com", to, msg)
	if err != nil {
		log.Fatal(err)
	} else {

		fmt.Println("Email enviado con exito...!!!!!")
	}

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
	fmt.Println("BACKEND GEOCHAT....!!!")
	//w.Write([]byte(`{"BACKEND GEOCHAT....!!!"}`))

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
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	pool := websocket.NewPool()
	go pool.Start()
	serveWs(pool, w, r)
}

// ****************     Definiendo rutas    ************************
func (ur *UserRouter) Routes() http.Handler {
	r := chi.NewRouter()

	// Configurar el middleware CORS para permitir todas las solicitudes desde cualquier origen

	r.Post("/login", ur.UserLogin)
	r.Post("/", ur.UserSignup) // 5555/api/v2/users/
	r.Post("/api/user/mail", ur.UserMail)

	r.Get("/ws", websocketHandler)

	return r
}
