package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var (
	DB *sql.DB
)

// Подключаемся к базе данных
func Connect(connStr string) {
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Ошибка проверки соединения с базой данных: %v", err)
	}

	log.Println("Успешное подключение к базе данных!")
}

// InitTables инициализирует таблицы
func InitTables() {
	// Создание таблицы пользователей
	query := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        nickname VARCHAR(50) NOT NULL,
		email VARCHAR(100) NOT NULL,
        password VARCHAR(50) NOT NULL
    );
    CREATE TABLE IF NOT EXISTS posts (
        id SERIAL PRIMARY KEY,
        user_id INTEGER REFERENCES users(id),
		date DATE NOT NULL,
        content TEXT NOT NULL,
		likes INTEGER NOT NULL,
		comments INTEGER NOT NULL
    );
    CREATE TABLE IF NOT EXISTS comments (
        id SERIAL PRIMARY KEY,
        user_id INTEGER REFERENCES users(id),
        post_id INTEGER REFERENCES posts(id),
		date DATE NOT NULL,
        content TEXT NOT NULL,
		likes INTEGER NOT NULL,
    );`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Ошибка инициализации таблиц:", err)
	}
}
