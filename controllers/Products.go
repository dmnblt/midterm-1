package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Products struct {
	ID          int64   `json:"id"`
	UserId      int64   `json:"user_id"`
	ProductName string  `json:"product_name"`
	Rating      int64   `json:"rating"`
	Price       float64 `json:"price"`
	Available   bool    `json:"available"`
	CreatedAt   string  `json:"created_at"`
}

func AddProduct(c *gin.Context) {
	var product Products
	err := c.ShouldBindJSON(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := dbConnect().Exec("INSERT INTO products (user_id, product_name, rating, price, available) VALUES (?, ?, ?, ?, ?)",
		product.UserId, product.ProductName, product.Rating, product.Price, product.Available)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error while inserting": err.Error()})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	product.ID = id
	c.JSON(http.StatusOK, product)
}

func GetProductsByUserId(c *gin.Context) {
	searchQuery := c.Query("userId")
	if searchQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing search query parameter"})
		return
	}

	rows, err := dbConnect().Query("SELECT id, user_id, product_name, rating, price, available, created_at FROM products WHERE user_id = ? ORDER BY rating DESC", searchQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	products := make([]Products, 0)

	for rows.Next() {
		var product Products
		err := rows.Scan(&product.ID, &product.UserId, &product.ProductName, &product.Rating, &product.Price, &product.Available, &product.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		products = append(products, product)
	}

	err = rows.Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func UpdateProductOwner(c *gin.Context) {
	userId := c.Query("userId")
	id := c.Query("id")
	if userId == "" || id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing search query parameter"})
		return
	}

	query := "UPDATE products SET user_id = ? WHERE id = ?"
	_, err := db.Exec(query, userId, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user email",
		})
		return
	}
	products := make([]Products, 0)

	c.JSON(http.StatusOK, products)
}

func GetProductsBetweenPrices(c *gin.Context) {
	from, err := strconv.ParseFloat(c.Query("from"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'from' parameter"})
		return
	}

	to, err := strconv.ParseFloat(c.Query("to"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'to' parameter"})
		return
	}

	rows, err := dbConnect().Query(`
		SELECT id, user_id, product_name, rating, price, available, created_at
		FROM products
		WHERE price BETWEEN ? AND ?
	`, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var products []Products
	for rows.Next() {
		var p Products
		err := rows.Scan(&p.ID, &p.UserId, &p.ProductName, &p.Rating, &p.Price, &p.Available, &p.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}
