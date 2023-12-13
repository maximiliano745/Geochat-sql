package v2

import (
	//"context"

	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	//"log"
	"net/http"
	"net/smtp"

	"github.com/github.com/maximiliano745/Geochat-sql/pkg/response"
	"github.com/github.com/maximiliano745/Geochat-sql/pkg/user"
	"github.com/github.com/maximiliano745/Geochat-sql/pkg/websocket"
	"github.com/go-chi/chi"
	"golang.org/x/crypto/bcrypt"
	//"github.com/maximiliano745/Geochat-sql/internal/server/v2/serverrtc"
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

func (ur *UserRouter) HaciendoTarea() {
	for {
		//fmt.Println("Realizando tarea: Fijandome los Pedidos de Amistad...")
		ur.Repository.ConsultaPedidosContacto()
		time.Sleep(3 * time.Second)
	}
}

func (ur *UserRouter) VerContactos(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("**** ACA Trayendo los Nombres de los Contactos...   ******")
	var u user.User
	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println("ids..: ", u.ID)
	defer r.Body.Close()
	ctx := r.Context()

	contactos, err := ur.Repository.GetOne(ctx, u.ID)
	if err != nil {
		fmt.Println(err)
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println("Contactos Trayendo: ", contactos.Username)

	// Ahora, escribimos el Username en la respuesta HTTP
	w.Header().Set("Content-Type", "application/json")
	responseJSON := map[string]string{"Username": contactos.Username}
	json.NewEncoder(w).Encode(responseJSON)
}

func (ur *UserRouter) UserContactos(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("**** ACA Viendo Los Concatsos...   ******")

	var u user.User
	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println("Contactos del id: ", u.ID)

	defer r.Body.Close()

	ctx := r.Context()

	// Llama a GetContactos para obtener los contactos del usuario
	contactos, err := ur.Repository.GetContactos(ctx, u.ID)
	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Print("Contactos Viendo: ", contactos)

	// Envía los contactos como respuesta en formato JSON como un arreglo
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contactos)
}

func (ur *UserRouter) UserMail(w http.ResponseWriter, r *http.Request) {

	var request struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Message string `json:"message"`
		Otro    string `json:"Otro"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	fmt.Println("\n\n\n ")
	fmt.Println("**** ACA ESTAMOS EN EL EMAIL   ******")

	fmt.Println("------------->Mail del que manda---> ", request.Email, "Recibe: ", request.Email)

	auth := smtp.PlainAuth("", "maxiargento745@gmail.com", "rwkycxemzftxidxi", "smtp.gmail.com")

	to := []string{request.Email}
	msg := []byte("Enviado por: " + request.Otro + "\r\n" +
		"Desde: Geochat..!!!\r\n" +
		"\r\n" +
		"Esto es la Invitacion de Contacto de GEOCHAT de: " + request.Otro + "---------------->   " + "https://maxi-geochat.onrender.com/" + "\n" + request.Message)
	err = smtp.SendMail("smtp.gmail.com:587", auth, "maxiargento745@gmail.com", to, msg)
	if err != nil {
		log.Fatal("Error:  --> ", err)
	} else {

		fmt.Println("Email enviado con exito...!!!!!")
		ctx := r.Context()

		// Obtener el primer usuario por correo electrónico
		userAcepta, err := ur.Repository.GetByMail(ctx, request.Email)
		if err != nil {
			// Manejar el error, por ejemplo, devolver un error HTTP o registrar un error
			return
		}

		// Obtener el segundo usuario por correo electrónico
		userOfrece, err := ur.Repository.GetByMail(ctx, request.Otro)
		if err != nil {
			// Manejar el error
			return
		}
		fmt.Println("ID Manda:", userOfrece.ID, "  ID Recibe:", userAcepta.ID)
		err = ur.Repository.AgregaPedidoAmistad(ctx, userOfrece.ID, userAcepta.ID)
		if err != nil {
			fmt.Println("Error: ", err)
		}
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
	fmt.Println("\n\n\n ")
	fmt.Println("**** ACA ESTAMOS EN EL LOGIN   ******")

	ctx := r.Context()
	uu, err = ur.Repository.GetByMail(ctx, uu.Email)

	if err == nil {
		response := map[string]interface{}{
			"message": "OK!!, Email EXISTENTE....!!!",
			"name":    uu.Username,
			"id":      uu.ID,
		}
		fmt.Println(" OK!!, Email EXISTENTE....!!!")
		resp, err := json.Marshal(response)
		if err != nil {
			// Manejar el error de conversión JSON
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("Nombre:", uu.Username, "     ID:", uu.ID, "      Password:", uu.Password, "    Email:", uu.Email)
		fmt.Println("Revisando Contraseña.....", u.Password, "     ", uu.Hash)

		if uu.PasswordMatch(u.Password) {
			fmt.Println("Contraseña correcta")
			w.Write(resp)
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
	fmt.Println("\n\n\n ")
	fmt.Println("\n *** ACA ESTAMOS EN EL REGISTRO DE USUAROIOS ****\nn")
	var u, uu user.User
	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		fmt.Print(err)
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	uu = u
	defer r.Body.Close()

	ctx := r.Context()
	uu, err = ur.Repository.GetByMail(ctx, uu.Email)
	if err == nil {
		fmt.Println("\nError Email EXISTENTE....!!!")
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
			fmt.Println("Nombre:", u.Username, "     ID:", u.ID, "      Password:", u.Password, "    Email:", u.Email)
			w.Write([]byte(`{"USUARIO CREADO CON  EXITO...!!!!"}`))
			w.WriteHeader(http.StatusOK)
		}
	}
}

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n\n\n ")
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

func websocketHandler(w http.ResponseWriter, r *http.Request, pool *websocket.Pool) {
	//pool := websocket.NewPool()
	//go pool.Start()
	serveWs(pool, w, r)
}

func (ur *UserRouter) CrearGrupos(w http.ResponseWriter, r *http.Request) {
	fmt.Print("\n**** ACA Creando Grupo y Guardando Integrantes   ******\n")

	var g user.Grupo
	err := json.NewDecoder(r.Body).Decode(&g)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println("datos: ", g)

	//Inserta el grupo y obtiene su ID
	grupoID, err := ur.Repository.CrGrupo(r.Context(), g)
	if err != nil {
		fmt.Println("Error:", err)
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
		return
	} else {
		fmt.Print("\nGuardado del GRUPO Exitoso....\n")
	}

	for _, contacto := range g.Contactos {
		fmt.Printf("ID: %d, Nombre: %s\n", contacto.ID, contacto.Nombre)
	}

	fmt.Print("\nDatos recibidos:\n")
	fmt.Println("Nombre del grupo:", g.Nombre)
	fmt.Println("Id del dueño:", g.IDueño)
	fmt.Println("ID del grupo:", grupoID)

	fmt.Println("")

	// Si todo está bien, puedes responder con un mensaje de éxito
	response.JSON(w, r, http.StatusOK, "Grupo creado exitosamente")
}

func (ur *UserRouter) VerGrupos(w http.ResponseWriter, r *http.Request) {
	fmt.Print("\nVer Grupos----------------------------------------------->")
	var u user.User
	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println("Grupos del id: ", u.ID)

	defer r.Body.Close()
	ctx := r.Context()

	grupos, err := ur.Repository.TraeGrupos(ctx, u.ID)
	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Print("\nGrupos Viendo: ", grupos, "\n\n")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grupos)
}

func (ur *UserRouter) TraeMiembrosGrupo(w http.ResponseWriter, r *http.Request) {

	fmt.Print("\nTrae Miembros de los Grupos----------------------------------------------->")
	var u user.User
	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println("************** Miembros de los Grupos del id:************* ", u.ID)
	defer r.Body.Close()
	ctx := r.Context()

	miembros, err := ur.Repository.TraeGruposMiembros(ctx, u.ID)
	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(miembros)

}

// ****************     Definiendo rutas    ************************
func (ur *UserRouter) Routes() http.Handler {
	r := chi.NewRouter()

	//go ur.HaciendoTarea()
	// Crea un WaitGroup para esperar a que ambas tareas terminen antes de salir
	var wg sync.WaitGroup

	// Inicia la tarea HaciendoTarea en segundo plano
	wg.Add(1)
	go func() {
		defer wg.Done()
		ur.HaciendoTarea()
	}()

	//Inicia la tarea ConsultaPedidosContacto en segundo plano
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	ur.UserRepository.ConsultaPedidosContacto()
	// }()

	// Roles

	r.Post("/login", ur.UserLogin)
	r.Post("/register", ur.UserSignup) // /api/v2/users/
	r.Post("/api/user/mail", ur.UserMail)
	r.Post("/contactos", ur.UserContactos)
	r.Post("/verContactos", ur.VerContactos)
	r.Post("/crearGrupos", ur.CrearGrupos)
	r.Post("/vergrupos", ur.VerGrupos)
	r.Post("/traerMiembrosGrupo", ur.TraeMiembrosGrupo)

	// r.Post("/create", serverrtc.CreateRoomRequestHandle)
	// r.Post("/join", serverrtc.JoinRoomRequestHandle)

	pool := websocket.NewPool()
	go pool.Start()

	r.Get("/wss", func(w http.ResponseWriter, r *http.Request) {
		websocketHandler(w, r, pool)
	})

	//r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
	//	websocketHandler(w, r, pool)
	//})

	return r
}
