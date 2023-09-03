package main

import (
	"os"
	"os/signal"

	"log"

	"github.com/github.com/maximiliano745/Geochat-sql/internal/data"

	"github.com/github.com/maximiliano745/Geochat-sql/internal/server"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	port := os.Getenv("PORT")
	serv, err := server.New(port)
	//serv, err := server.New("8000")

	if err != nil {
		log.Fatal(err)
	}

	// Coneccion a la Base de Datos.
	d := data.New()
	if err := d.DB.Ping(); err != nil {
		log.Fatal(err)
	}

	// inicia el servidor.En Una 'Go Rutina'
	go serv.Start()

	// Espere una interrupción. Espera de recibir un evento del sistema (Ej: ctrl+C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Intenta un cierre elegante. // Podemos hacer mas cosas Cuando Sale....(Ej: ctrl+C)
	serv.Close()
	//data.Close()

}
