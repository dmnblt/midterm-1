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
