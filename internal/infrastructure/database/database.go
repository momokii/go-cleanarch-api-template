package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	var err error

	host := os.Getenv("HOST_POSTGRES")
	port := os.Getenv("PORT_POSTGRES")
	user := os.Getenv("USER_POSTGRES")
	password := os.Getenv("PASSWORD_POSTGRES")
	dbname := os.Getenv("DATABASE_POSTGRES")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, port, user, password, dbname)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Connect postgres error: ", err)
	}

	DB.SetMaxIdleConns(10)
	DB.SetMaxOpenConns(10)

	err = DB.Ping()
	if err != nil {
		log.Fatal("Ping postgres error: ", err)
	}

	fmt.Println("Successfully connected to database")
}
