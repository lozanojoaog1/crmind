package database

import (
	"database/sql"
	"log"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	connStr := "host=p-2hrz9m0lvu.pg.biganimal.io port=5432 user=edb_admin password=CRMind*2024*898 dbname=edb_admin sslmode=require"
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("Falha ao conectar ao banco de dados:", err)
	}
	log.Println("Conexão com o banco de dados estabelecida com sucesso")
}
