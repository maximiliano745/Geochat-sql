package v1

import (
	"net/http"

	"github.com/github.com/maximiliano745/Geochat-sql/internal/data"
	"github.com/go-chi/chi"
)

func New() http.Handler {
	r := chi.NewRouter()

	ur := &UserRouter{
		Repository: &data.UserRepository{
			Data: data.New(),
		},
	}
	r.Mount("/Maxi", ur.Routes()) //  http://localhost:10000/api/v1/Maxi/

	/* las rutas que definimos dentro de Routes ahora tendrán como base /users , por ejemplo, para el caso de
	r.Get("/", ur.GetAllHandler), la ruta será /Maxi/
	*/

	return r
}
