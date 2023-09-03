package data

import (
	"database/sql"
	"log"
	"sync"
)

var (
	data *Data
	once sync.Once
)

// Los datos gestionan la conexión a la base de datos.
type Data struct {
	DB *sql.DB
}

func initDB() {
	db, err := getConnection()
	if err != nil {
		log.Panic(err)
	}

	err = MakeMigration(db)
	if err != nil {
		log.Panic(err)
	}

	data = &Data{
		DB: db,
	}
}

func New() *Data {
	once.Do(initDB)

	return data
}
