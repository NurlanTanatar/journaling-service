package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func localInitDb() string {

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_DBNAME")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

func openDB() error {
	psqlInfo := localInitDb()
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}
	DB = db
	fmt.Println("Successfully connected to db")
	return nil
}

func closeDB() error {
	err := DB.Close()
	if err != nil {
		return err
	}
	fmt.Println("Connection to db is closed")
	return nil
}

func setupDB() error {
	_, err := DB.Exec(`create table if not exists journal(id serial primary key, 
														  title text, 
														  completed boolean default false, 
														  position integer, 
														  criticality integer,
														  dateStart timestamp,
														  dateEnd timestamp
					 									);`)

	if err != nil {
		return err
	}
	return nil
}
