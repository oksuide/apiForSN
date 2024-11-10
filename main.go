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
	// Подключаемся к .env
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

	// Группа маршрутов для авторизованных пользователей
	authorized := router.Group("/api")
	authorized.Use(middleware.AuthMiddleware())

	// Группа маршрутов для пользователей
	users := authorized.Group("/users")
	{
		users.GET("/:id", handlers.GetUser)       // Получение пользователя по ID
		users.PUT("/:id", handlers.UpdateUser)    // Обновление данных пользователя по ID
		users.DELETE("/:id", handlers.DeleteUser) // Удаление пользователя по ID
		users.POST("/", handlers.CreateUser)      // Создание пользователя
	}

	// Группа маршрутов для постов
	posts := authorized.Group("/posts")
	{
		posts.POST("/", handlers.CreatePost)           // Создание поста
		posts.GET("/:postID", handlers.GetPost)        // Получение поста по ID
		posts.PUT("/:postID", handlers.UpdatePost)     // Обновление поста по ID
		posts.DELETE("/:postID", handlers.DeletePost)  // Удаление поста по ID
		posts.POST("/:postID/like", handlers.LikePost) // Лайк поста
	}

	// Группа маршрутов для комментариев
	comments := authorized.Group("/comments")
	{
		comments.POST("/posts/:postID/comments/", handlers.CreateComment)              // Создание комментария
		comments.GET("/posts/:postID/comments/:commentID", handlers.GetComment)        // Получение комментария по ID
		comments.PUT("/posts/:postID/comments/:commentID", handlers.UpdateComment)     // Обновление комментария по ID
		comments.DELETE("/posts/:postID/comments/:commentID", handlers.DeleteComment)  // Удаление комментария по ID
		comments.POST("/posts/:postID/comments/:commentID/like", handlers.LikeComment) // Лайк комментария
	}

	// Запуск сервера на порту 8080
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
