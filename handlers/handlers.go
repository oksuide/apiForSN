package handlers

import (
	"apiForSN/db"
	"apiForSN/models"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

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
	// Логика для получения пользователя
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

func UpdateUser(c *gin.Context) {
	// Логика для обновления пользователя
	// Получаем ID пользователя из контекста
	userID := c.MustGet("userID").(int)

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
	userID := c.MustGet("userID").(int)
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
