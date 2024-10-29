package main

import (
	"apiForSN/handlers"
	"apiForSN/middleware"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/oksuide/apiForSN/db"
)

func main() {
	// Подключаемся в дотенву
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Создание роутера
	router := gin.Default()
	// Подключение к базе данных
	connStr := os.Getenv("ConnStr")

	db.Connect(connStr)
	db.InitTables()

	// Применяем AuthMiddleware ко всем маршрутам, требующим авторизации
	authorized := router.Group("/api")
	authorized.Use(middleware.AuthMiddleware())
	{
		// Чтение данных о пользователе по ID
		authorized.GET("/user", handlers.GetUser)
		// Обновление данных пользователя по ID
		authorized.PUT("/user", handlers.UpdateUser)
		// Удаление пользователя по IDhandlershandlers
		authorized.DELETE("/users/:id", handlers.DeleteUser)
		// Создание пользователя
		authorized.POST("/users", handlers.CreateUser)
	}

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
