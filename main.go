package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Allanc2002/P_Like_publication/connection"
	"github.com/Allanc2002/P_Like_publication/functions"
)

var SECRET_KEY string

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found")
	}

	SECRET_KEY = os.Getenv("SECRET_KEY")
	if SECRET_KEY == "" {
		log.Fatal("SECRET_KEY not set in environment")
	}

	db, err := connection.MongoConnection()
	if err != nil {
		log.Fatal("Could not connect to DB:", err)
	}

	r := gin.Default()

	// Endpoint
	r.POST("/like", func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token missing or invalid"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
			return
		}
		userID := int(userIDFloat)

		var json struct {
			IdPublication string `json:"_id"`
		}

		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		message, code, err := functions.LikePublication(db, json.IdPublication, userID)
		if err != nil {
			c.JSON(code, gin.H{"error": err.Error()})
			return
		}

		c.JSON(code, gin.H{"message": message})
	})

	r.Run(":8080")
}
