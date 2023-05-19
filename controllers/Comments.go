package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Comments struct {
	ID         int    `json:"id"`
	Product_ID int    `json:"product_id"`
	User_ID    int64  `json:"user_id"`
	Title      string `json:"title"`
	Body       string `json:"body"`
	Likes      int64  `json:"comment_likes"`
	CreatedAt  string `json:"created_at"`
}

func AddComment(c *gin.Context) {
	var comment Comments
	err := c.ShouldBindJSON(&comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := dbConnect().Exec("INSERT INTO comments (user_id, product_id, title, body) VALUES (?, ?, ?, ?)",
		comment.User_ID, comment.Product_ID, comment.Title, comment.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error while inserting": err.Error()})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	comment.ID = int(id)
	c.JSON(http.StatusOK, comment)
}

func GetCommentsByProductID(c *gin.Context) {
	searchQuery := c.Query("product_id")
	if searchQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing search query parameter"})
		return
	}

	rows, err := dbConnect().Query("SELECT id, product_id, user_id ,title, body, comment_likes ,created_at FROM comments WHERE product_id = ?", searchQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	comments := make([]Comments, 0)

	for rows.Next() {
		var comment Comments
		err := rows.Scan(&comment.ID, &comment.Product_ID, &comment.User_ID, &comment.Title, &comment.Body, &comment.Likes, &comment.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		comments = append(comments, comment)
	}

	err = rows.Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}

func GetCommentsByLikesDesc(c *gin.Context) {

	rows, err := dbConnect().Query(`
		SELECT *
		FROM comments
		ORDER BY comment_likes DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var comments []Comments
	for rows.Next() {
		var comment Comments
		err := rows.Scan(&comment.ID, &comment.Product_ID, &comment.User_ID, &comment.Title, &comment.Body, &comment.Likes, &comment.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}

func LikeComment(c *gin.Context) {

	searchComment := c.Query("id")
	if searchComment == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing search query parameter"})
		return
	}
	rows, err := dbConnect().Query("SELECT comment_likes FROM comments WHERE id=?", searchComment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
		return
	}

	defer rows.Close()
	comment_likes := 0
	for rows.Next() {
		err := rows.Scan(&comment_likes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
			return
		}

	}
	comment_likes += 1

	result, err := dbConnect().Query("UPDATE comments SET comment_likes =?", comment_likes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error3": err.Error()})
		return
	}
	result.Close()
	c.JSON(http.StatusOK, searchComment)
}
