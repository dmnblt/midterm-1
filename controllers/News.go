package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type News struct {
	ID      int    `json:"id"`
	IconURL string `json:"icon_url"`
	Title   string `json:"title"`
	Body    string `json:"body"`
}

func AddNews(c *gin.Context) {
	isAdmin := isAdmin()
	if !isAdmin {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
		return
	}

	var news News
	if err := c.ShouldBindJSON(&news); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	db := dbConnect()
	defer db.Close()

	_, err := db.Exec("INSERT INTO news (icon_url, title, body) VALUES (?, ?, ?)", news.IconURL, news.Title, news.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to insert news into database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "news added successfully"})
}
