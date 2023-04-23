# midterm-1
Go project for midterm exam

REPORT 8

Added new method for changing users password. By gin framework we get id of the user and by that id we will change it:
1.create new password 
2. convert it into hash 
3. upload it to out database

func UpdateUserPassword(c *gin.Context) {
	userID := c.Param("id")

	var user User
	query := "SELECT * FROM users WHERE id = ?"   err := db.QueryRow(query, userID).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",      })
		return   }

	var update struct {
		Password string `json:"password" binding:"required"`   }
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",      })
		return   }
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(update.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return   }
	query = "UPDATE users SET password = ? WHERE id = ?"   _, err = db.Exec(query, hashedPassword, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user email",      })
		return   }

	user.Password = hashedPassword
	c.JSON(http.StatusOK, user)
}

REPORT 9
Added a method and created a news table 
for admins that can post news and a new column 
in the users isAdmin() table.  And made this parameter 
to be determined when registering users

import (
"github.com/gin-gonic/gin"
"net/http"
)

type News struct {
ID       int    `json:"id"`
IconURL  string `json:"icon_url"`
Title    string `json:"title"`
Body     string `json:"body"`
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


REPORT 10
Added method for deleting and updating news.

func deleteNews(c *gin.Context) {
id := c.Param("id") // get the id of the news item to be deleted from the URL parameter
var news News
db := dbConnect() // get a database connection
if err := db.Where("id = ?", id).Delete(&news).Error; err != nil {
// handle error if any
c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}
c.JSON(http.StatusOK, gin.H{"message": "News item deleted successfully"})
}

func updateNews(c *gin.Context) {
id := c.Param("id")

	var news News
	db := dbConnect()

	// Find the news item in the database by its ID
	if err := db.First(&news, id).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "News item not found"})
		return
	}

	// Parse the JSON request body into a News struct
	var updatedNews News
	if err := c.ShouldBindJSON(&updatedNews); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	// Update the news item in the database
	news.IconURL = updatedNews.IconURL
	news.Title = updatedNews.Title
	news.Body = updatedNews.Body

	if err := db.Save(&news).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": news})
}