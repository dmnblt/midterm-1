package controllers

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type User struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

var db *sql.DB

func dbConnect() *sql.DB {
	var err error

	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/go_midterm")
	if err != nil {
		log.Fatalf("Error connecting to database: %v \n", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error checking the database connection: %v \n", err)
	}
	fmt.Println("Connected to database!")

	return db
}

func RegisterUser(c *gin.Context) {
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var count int
	err := dbConnect().QueryRow("SELECT COUNT(*) FROM users WHERE email=?", newUser.Email).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	_, err = dbConnect().Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", newUser.Name, newUser.Email, hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting user into database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func AuthorizeUser(c *gin.Context) {
	defer db.Close()
	var authUser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&authUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dbPassword string
	err := dbConnect().QueryRow("SELECT password FROM users WHERE email=?", authUser.Email).Scan(&dbPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(authUser.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Authorized successfully"})
}

func SearchUsersByName(c *gin.Context) {
	defer db.Close()
	searchQuery := c.Query("name")
	if searchQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing search query parameter"})
		return
	}

	rows, err := dbConnect().Query("SELECT id, name, email, created_at FROM users WHERE name LIKE ?", "%"+searchQuery+"%")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying database"})
		return
	}

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning row"})
			return
		}
		users = append(users, user)

		fmt.Println(users)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func FindUsersBetweenDates(c *gin.Context) {
	var users []User

	start := c.Query("start_date")
	end := c.Query("end_date")

	query := "SELECT * FROM users WHERE created_at BETWEEN ? AND ?"
	rows, err := db.Query(query, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get users",
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get users",
			})
			return
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}

func UpdateUserEmail(c *gin.Context) {
	userID := c.Param("id")

	var user User
	query := "SELECT * FROM users WHERE id = ?"
	err := db.QueryRow(query, userID).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	var update struct {
		Email string `json:"email" binding:"required"`
	}
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

	query = "UPDATE users SET email = ? WHERE id = ?"
	_, err = db.Exec(query, update.Email, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user email",
		})
		return
	}

	user.Email = update.Email
	c.JSON(http.StatusOK, user)
}
func UpdateUserPassword(c *gin.Context) {
	userID := c.Param("id")

	var user User
	var query = "SELECT * FROM users WHERE id = ?"
	err := db.QueryRow(query, userID).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found"})
		return
	}

	var update struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload"})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(update.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}
	query = "UPDATE users SET password = ? WHERE id = ?"
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user email"})
		return
	}

	user.Password = string(hashedPassword)
	c.JSON(http.StatusOK, user)
}
