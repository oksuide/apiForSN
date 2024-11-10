package main

import (
	"apiForSN/db"
	"apiForSN/handlers"
	"apiForSN/middleware"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Подключаемся к .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Подключение к базе данных
	connStr := os.Getenv("ConnStr")
	db.Connect(connStr)
	db.InitTables()

	// Создание роутера
	router := gin.Default()

	// Применяем AuthMiddleware ко всем маршрутам, требующим авторизации
	authorized := router.Group("/api")
	authorized.Use(middleware.AuthMiddleware())
	{
		// Роуты для работы с пользователями
		authorized.GET("/user", handlers.GetUser)
		authorized.PUT("/user", handlers.UpdateUser)
		authorized.DELETE("/users/:id", handlers.DeleteUser)
		authorized.POST("/users", handlers.CreateUser)

		// Роуты для постов
		posts := authorized.Group("/posts")
		{
			posts.Use(middleware.PostIDMiddleware()) // Применяем middleware для postID
			posts.GET("/:postID", handlers.GetPost)
			posts.PUT("/:postID", handlers.UpdatePost)
			posts.DELETE("/:postID", handlers.DeletePost)
			posts.POST("/", handlers.CreatePost)
			posts.POST("/:postID/like", handlers.LikePost)
		}

		// Роуты для комментариев
		comments := authorized.Group("/comments")
		{
			comments.Use(middleware.CommentIDMiddleware()) // Применяем middleware для commentID
			comments.GET("/:commentID", handlers.GetComment)
			comments.PUT("/:commentID", handlers.UpdateComment)
			comments.DELETE("/:commentID", handlers.DeleteComment)
			comments.POST("/", handlers.CreateComment)
			comments.POST("/:commentID/like", handlers.LikeComment)
		}
	}

	// Запуск сервера на порту 8080
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
