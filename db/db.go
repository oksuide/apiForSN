package db

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

// Подключаемся к базе данных
func Connect(connStr string) {
	var err error
	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}

	// Проверяем соединение с базой данных
	sqlDB, err := DB.DB()
	if err != nil {
		fmt.Println("Error getting DB instance:", err)
		return
	}

	if err := sqlDB.Ping(); err != nil {
		fmt.Println("Database connection failed:", err)
	} else {
		fmt.Println("Database connection successful!")
	}

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
		likes INTEGER NOT NULL
    );
	CREATE TABLE IF NOT EXISTS likes (
        id SERIAL PRIMARY KEY,
        user_id INTEGER REFERENCES users(id),
        post_id INTEGER REFERENCES posts(id),
		comment_id INTEGER REFERENCES comments(id)
	);`
	err := DB.Exec(query).Error
	if err != nil {
		log.Fatal("Ошибка инициализации таблиц:", err)
	}
}
