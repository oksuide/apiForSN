package main

import (
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"github.com/oksuide/apiForSN/db"
)

func main() {
	router := gin.Default()
	// Подключение к базе данных
	connStr := "user=username password=yourpassword dbname=yourdb sslmode=disable"
	db.Connect(connStr)
	db.InitTables()

	// Запуск сервера на порту 8080
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

// // Простой пример маршрута, где /ping - путь, 200 - HTTP статус, что все ок
// // Сокращение для map[string]interface{}, которое используется для создания JSON-ответа.
// // Создаем JSON с полем "message" и значением "pong".

// router.GET("/ping", func(c *gin.Context) {
//     c.JSON(200, gin.H{
//         "message": "pong",
//     })
// })
