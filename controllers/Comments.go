package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Comments struct {
	ID        int    `json:"id"`
	News_ID   int    `json:"news_id"`
	User_ID   int64  `json:"user_id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
}

func AddComment(c *gin.Context) {
	var comment Comments
	err := c.ShouldBindJSON(&comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := dbConnect().Exec("INSERT INTO comments (user_id, news_id, title, body) VALUES (?, ?, ?, ?)",
		comment.User_ID, comment.News_ID, comment.Title, comment.Body)
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

func GetCommentsByNewsID(c *gin.Context) {
	searchQuery := c.Query("news_id")
	if searchQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing search query parameter"})
		return
	}

	rows, err := dbConnect().Query("SELECT id, news_id, user_id ,title, body, created_at FROM comments WHERE news_id = ?", searchQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	comments := make([]Comments, 0)

	for rows.Next() {
		var comment Comments
		err := rows.Scan(&comment.ID, &comment.News_ID, &comment.User_ID, &comment.Title, &comment.Body, &comment.CreatedAt)
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
