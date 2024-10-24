package models

type User struct {
	ID       int    `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Post struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id"`
	Date     int    `json:"date"`
	Content  string `json:"content"`
	Likes    int    `json:"likes"`
	Comments int    `json:"comments"`
}

type Comment struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_id"`
	PostID  int    `json:"post_id"`
	Date    int    `json:"date"`
	Content string `json:"content"`
	Likes   int    `json:"likes"`
}
