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
	r.Mount("/users", ur.Routes()) //  http://localhost:9000/api/v1/users/

	/* las rutas que definimos dentro de Routes ahora tendrán como base /users , por ejemplo, para el caso de
	r.Get("/", ur.GetAllHandler), la ruta será /users/
	*/

	/* pr := &PostRouter{
	       Repository: &data.PostRepository{
	           Data: data.New(),
	       },
	   }

	r.Mount("/posts", pr.Routes()) */

	return r
}
