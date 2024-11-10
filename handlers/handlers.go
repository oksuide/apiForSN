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
func CreatePost(c *gin.Context) {
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
	postID, _ := c.Get("postID")
	userID, _ := c.Get("userID")
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
		// Удаляем пост
		if err := db.DB.Delete(&post).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to delete post"})
			return
		}

		// Удаляем все комментарии, связанные с постом
		if err := db.DB.Where("post_id = ?", postID).Delete(&models.Comment{}).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to delete comments"})
			return
		}

		c.JSON(200, gin.H{"message": "Post and related comments deleted successfully"})
	} else {
		c.JSON(403, gin.H{"error": "You must be the author of the post"})
	}
}

func UpdatePost(c *gin.Context) {
	// Получаем ID поста и ID пользователя из контекста
	postID, _ := c.Get("postID")
	userID, _ := c.Get("userID")

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

func GetPost(c *gin.Context) {
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
		"id":      existingPost.ID,
		"content": existingPost.Content,
		"userID":  existingPost.UserID,
	})
}

// Блок работы с комментариями
func CreateComment(c *gin.Context) {
	var comment models.Comment
	// Привязываем JSON-данные из тела запроса к переменной comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Устанавливем текущие дата/время для комментария
	comment.Date = int(time.Now().Unix())
	// Устанавливаем userID и postID из контекста
	postID, postExists := c.Get("postID")
	userID, userExists := c.Get("userID")
	if !postExists || !userExists {
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}
	comment.UserID = userID.(int)
	comment.PostID = postID.(int)
	// Сохраняем комментарий в базе данных
	if err := db.DB.Create(&comment).Error; err != nil {
		c.JSON(500, gin.H{"error": "Error creating comment"})
		return
	}
	// Увеличиваем количество комментариев в посте на 1
	if err := db.DB.Model(&models.Post{}).Where("id = ?", comment.PostID).UpdateColumn("comments", gorm.Expr("comments + ?", 1)).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update post comments count"})
		return
	}

	// Возвращаем успешный ответ с данными о созданном посте
	c.JSON(201, gin.H{
		"id":      comment.ID,
		"user_id": comment.UserID,
		"post_id": comment.PostID,
		"date":    comment.Date,
		"content": comment.Content,
	})
}

func DeleteComment(c *gin.Context) {
	commentID, _ := c.Get("commentID")
	userID, _ := c.Get("userID")

	// Ищем комментарий в базе данных по ID
	var comment models.Comment
	if err := db.DB.Where("id = ?", commentID).First(&comment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//
			c.JSON(404, gin.H{"error": "Comment not found"})
			return
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}
	// Проверяем автора и удаляем комментарий
	if userID == comment.UserID {
		if err := db.DB.Delete(&comment).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to delete comment"})
			return
		}
		c.JSON(200, gin.H{"message": "Comment deleted successfully"})
	} else {
		c.JSON(403, gin.H{"error": "You must be author of the comment"})
		return
	}
	// Уменьшаем количество комментариев в посте на 1
	if err := db.DB.Model(&models.Post{}).Where("id = ?", comment.PostID).UpdateColumn("comments", gorm.Expr("comments - ?", 1)).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update post comments count"})
		return
	}
}

func UpdateComment(c *gin.Context) {
	// Получаем ID комментария и ID пользователя из контекста
	commentID, _ := c.Get("commentID")
	userID, _ := c.Get("userID")

	// Ищем комментарий в базе данных по ID
	var comment models.Post
	if err := db.DB.Where("id = ?", commentID).First(&comment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "Comment not found"})
			return
		}
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}

	// Проверяем, что пользователь является автором комментария
	if userID != comment.UserID {
		c.JSON(403, gin.H{"error": "You must be the author of the comment to update it"})
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

	// Проверяем, что новое содержание комментария не пустое
	if updateData.Content == "" {
		c.JSON(400, gin.H{"error": "New comment content can't be empty"})
		return
	}

	// Обновляем контент и дату поста
	updates := map[string]interface{}{
		"content": updateData.Content,
		"date":    int(time.Now().Unix()),
	}
	if err := db.DB.Model(&comment).Updates(updates).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update comment"})
		return
	}

	// Обновляем значение `comment` для возврата обновленных данных
	comment.Content = updateData.Content
	comment.Date = updates["date"].(int)

	// Возвращаем успешный ответ с обновленными данными поста
	c.JSON(200, gin.H{
		"content": comment.Content,
		"date":    comment.Date,
	})
}

func GetComment(c *gin.Context) {
	// Получаем ID пользователя из параметров URL
	commentID, err := strconv.Atoi(c.Param("commentID"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid comment ID"})
		return
	}
	// Проверяем наличие комментария с таким ID
	var existingComment models.Comment
	if err := db.DB.Where("id = ?", commentID).First(&existingComment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "Post not found"})
			return
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(200, gin.H{
		"id":      existingComment.ID,
		"userID":  existingComment.UserID,
		"postID":  existingComment.PostID,
		"content": existingComment.Content,
	})
}

// Блок работы с лайками
func LikePost(c *gin.Context) {
	userID, _ := c.Get("userID")
	postID, _ := c.Get("postID")

	// Проверка, ставил ли уже лайк этот пользователь, и удаление лайка, если он стоит
	var like models.Like
	if err := db.DB.Where("user_id = ? AND post_id = ?", userID, postID).First(&like).Error; err == nil {
		// Лайк найден, значит удаляем его
		if err := db.DB.Delete(&like).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to remove like"})
			return
		}

		// Уменьшаем количество лайков у поста
		if err := db.DB.Model(&models.Post{}).Where("id = ?", postID).Update("likes", gorm.Expr("likes - ?", 1)).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to update like count"})
			return
		}

		c.JSON(200, gin.H{"message": "Like removed successfully"})
		return
	}

	// Добавляем лайк, так как его еще нет
	newLike := models.Like{
		UserID: userID.(int),
		PostID: postID.(*int),
	}
	if err := db.DB.Create(&newLike).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to add like"})
		return
	}

	// Увеличиваем счётчик лайков у поста
	if err := db.DB.Model(&models.Post{}).Where("id = ?", postID).Update("likes", gorm.Expr("likes + ?", 1)).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update like count"})
		return
	}

	c.JSON(200, gin.H{"message": "Post liked successfully"})
}

func LikeComment(c *gin.Context) {
	userID, _ := c.Get("userID")
	commentID, _ := c.Get("commentID")

	// Проверка, ставил ли уже лайк этот пользователь и удалить в случае, если он стоит
	var like models.Like
	if err := db.DB.Where("user_id = ? AND comment_id = ?", userID, commentID).First(&like).Error; err == nil { // Здесь исправлено
		if err := db.DB.Delete(&like).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to remove like"})
			return
		}
		// Уменьшаем количество лайков в комментарии
		if err := db.DB.Model(&models.Comment{}).Where("id = ?", commentID).Update("likes", gorm.Expr("likes - ?", 1)).Error; err != nil { // Здесь исправлено
			c.JSON(500, gin.H{"error": "Failed to update like count"})
			return
		}
		c.JSON(200, gin.H{"message": "Like removed successfully"})
		return
	}

	// Добавляем лайк
	newLike := models.Like{
		UserID:    userID.(int),
		CommentID: commentID.(*int),
	}
	if err := db.DB.Create(&newLike).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to add like"})
		return
	}

	// Увеличиваем счётчик лайков комментария
	if err := db.DB.Model(&models.Comment{}).Where("id = ?", commentID).Update("likes", gorm.Expr("likes + ?", 1)).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update like count"})
		return
	}

	c.JSON(200, gin.H{"message": "Comment liked successfully"})
}
