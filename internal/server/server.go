// internal/server/server.go
package server

import (
	"log"
	"net/http"
	"time"

	v1 "github.com/github.com/maximiliano745/Geochat-sql/internal/server/v1"
	v2 "github.com/github.com/maximiliano745/Geochat-sql/internal/server/v2"
	"github.com/go-chi/chi"
	"github.com/rs/cors"
)

// Server is a base server configuration.
type Server struct {
	server *http.Server
}

// New inicialize a new server with configuration.
func New(port string) (*Server, error) {
	r := chi.NewRouter()

	c := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
		AllowedOrigins:   []string{"*"}, // All origins
	})

	r.Use(c.Handler) // Aplicar el middleware CORS al enrutador

	// API routes version 1.
	r.Mount("/api/v1", v1.New())

	// API routes version 1.
	r.Mount("/api/v2", v2.New())

	serv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server := Server{server: serv}

	return &server, nil
}

// Close server resources.
func (serv *Server) Close() error {
	// TODO: add resource closure.
	return nil
}

// Start the server.
func (serv *Server) Start() {

	log.Printf("Servidor ejecutándose en http://localhost%s", serv.server.Addr)
	log.Fatal(serv.server.ListenAndServe())
}
