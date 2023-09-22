package v3

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
	r.Mount("/users", ur.Routes()) //  http://localhost:5555/api/v3/users/

	return r
}
