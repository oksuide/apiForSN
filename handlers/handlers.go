package handlers

import (
	"apiForSN/db"
	"apiForSN/models"
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Блок работы с юзером
func CreateUser(c *gin.Context) {
	var user models.User
	// Привязываем JSON-данные из тела запроса к переменной user
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Проверяем уникальность email, чтобы не создать дублирующегося пользователя
	var existingUser models.User
	if err := db.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		c.JSON(400, gin.H{"error": "User with this email already exists"})
		return
	}

	// Хешируем пароль перед сохранением
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error hashing password"})
		return
	}
	user.Password = hashedPassword

	// Сохраняем пользователя в базе данных
	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Error creating user"})
		return
	}

	// Возвращаем успешный ответ с данными о созданном пользователе (без пароля)
	c.JSON(201, gin.H{
		"id":       user.ID,
		"nickname": user.Nickname,
		"email":    user.Email,
	})
}

func GetUser(c *gin.Context) {
	// Получаем ID пользователя из параметров URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}
	// Проверяем наличие пользователя с таким ID
	var existingUser models.User
	if err := db.DB.Where("id = ?", id).First(&existingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Пользователь с таким email не найден
			c.JSON(404, gin.H{"error": "User not found"})
			return
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(200, gin.H{
		"id":       existingUser.ID,
		"nickname": existingUser.Nickname,
		"email":    existingUser.Email,
	})
}

func newFunction(c *gin.Context) (models.User, bool) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return models.User{}, true
	}
	return user, false
}

func UpdateUser(c *gin.Context) {
	// Получаем ID пользователя из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	// Привязываем JSON с изменениями к структуре
	var updateData struct {
		Email    string `json:"email,omitempty"`
		Nickname string `json:"nickname,omitempty"`
		Password string `json:"password,omitempty"`
	}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	// Ищем пользователя в базе данных по ID
	var user models.User
	if err := db.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Пользователь с таким id не найден
			c.JSON(404, gin.H{"error": "User not found"})
			return
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}
	// Обновляем только те поля, которые присутствуют в запросе
	updates := make(map[string]interface{})
	if updateData.Email != "" {
		updates["email"] = updateData.Email
	}
	if updateData.Nickname != "" {
		updates["nickname"] = updateData.Nickname
	}
	if updateData.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateData.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to hash password"})
			return
		}
		updates["password"] = string(hashedPassword)
	}

	// Выполняем обновление в базе данных
	if len(updates) > 0 {
		if err := db.DB.Model(&user).Updates(updates).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to update user"})
			return
		}
	}

	// Возвращаем обновленную информацию о пользователе
	c.JSON(200, gin.H{
		"id":       user.ID,
		"nickname": user.Nickname,
		"email":    user.Email,
	})

}

func DeleteUser(c *gin.Context) {
	// Логика для удаления пользователя
	// Получаем ID пользователя из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}
	// Ищем пользователя в базе данных по ID
	var user models.User
	if err := db.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Пользователь с таким id не найден
			c.JSON(404, gin.H{"error": "User not found"})
			return
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}
	// Удаляем пользователя
	if err := db.DB.Delete(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(200, gin.H{"message": "User deleted successfully"})
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Блок работы с постами
func NewPost(c *gin.Context) {
	var post models.Post
	// Привязываем JSON-данные из тела запроса к переменной post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Устанавливем текущие дата/время для поста
	post.Date = int(time.Now().Unix())
	// Устанавливаем userID из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}
	post.UserID = userID.(int)

	// Сохраняем пост в базе данных
	if err := db.DB.Create(&post).Error; err != nil {
		c.JSON(500, gin.H{"error": "Error creating post"})
		return
	}

	// Возвращаем успешный ответ с данными о созданном посте
	c.JSON(201, gin.H{
		"id":      post.ID,
		"user_id": post.UserID,
		"date":    post.Date,
		"content": post.Content,
	})
}

func DeletePost(c *gin.Context) {
	// Получаем ID поста и ID пользователя из контекста
	postID, postExists := c.Get("postID")
	userID, userExists := c.Get("userID")
	if !postExists || !userExists {
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}
	// Ищем пост в базе данных по ID
	var post models.Post
	if err := db.DB.Where("id = ?", postID).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Пост с таким id не найден
			c.JSON(404, gin.H{"error": "Post not found"})
			return
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}
	// Проверяем автора и удаляем пост
	if userID == post.UserID {
		if err := db.DB.Delete(&post).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to delete user"})
			return
		}
		c.JSON(200, gin.H{"message": "Post deleted successfully"})
	} else {
		c.JSON(403, gin.H{"error": "You must be author of post"})
	}
}

func UpdatePost(c *gin.Context) {
	// Получаем ID поста и ID пользователя из контекста
	postID, postExists := c.Get("postID")
	userID, userExists := c.Get("userID")
	if !postExists || !userExists {
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}

	// Ищем пост в базе данных по ID
	var post models.Post
	if err := db.DB.Where("id = ?", postID).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}

	// Проверяем, что пользователь является автором поста
	if userID != post.UserID {
		c.JSON(403, gin.H{"error": "You must be the author of the post to update it"})
		return
	}

	// Привязываем JSON с изменениями к структуре
	var updateData struct {
		Content string `json:"content,omitempty"`
	}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	// Проверяем, что новое содержание поста не пустое
	if updateData.Content == "" {
		c.JSON(400, gin.H{"error": "New post content can't be empty"})
		return
	}

	// Обновляем контент и дату поста
	updates := map[string]interface{}{
		"content": updateData.Content,
		"date":    int(time.Now().Unix()),
	}
	if err := db.DB.Model(&post).Updates(updates).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update post"})
		return
	}

	// Обновляем значение `post` для возврата обновленных данных
	post.Content = updateData.Content
	post.Date = updates["date"].(int)

	// Возвращаем успешный ответ с обновленными данными поста
	c.JSON(200, gin.H{
		"content": post.Content,
		"date":    post.Date,
	})
}

func FindPost(c *gin.Context) {
	// Получаем ID пользователя из параметров URL
	postID, err := strconv.Atoi(c.Param("postID"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid post ID"})
		return
	}
	// Проверяем наличие пользователя с таким ID
	var existingPost models.Post
	if err := db.DB.Where("id = ?", postID).First(&existingPost).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "Post not found"})
			return
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(200, gin.H{
		"id":       existingPost.ID,
		"nickname": existingPost.Content,
		"userID":   existingPost.UserID,
	})
}
