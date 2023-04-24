package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Comments struct {
	ID      int    `json:"id"`
	News_ID int    `json:"news_id"`
	Title   string `json:"title"`
	Body    string `json:"body"`
}

func AddComment(c *gin.Context) {
	isAdmin := isAdmin()
	if !isAdmin {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
		return
	}

	var comment Comments
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	db := dbConnect()
	defer db.Close()
	var query = "SELECT * FROM news WHERE id = ?"
	if err := db.QueryRow(query, comment.News_ID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "News not found"})
		return
	}
	_, err := db.Exec("INSERT INTO comments (title, body) VALUES (?, ?, ?)", comment.Title, comment.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to insert comment into database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "comment added successfully"})
}
