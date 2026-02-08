package handlers

import (
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	PublicKey string    `json:"public_key"`
	CreatedAt time.Time `json:"created_at"`
}

type Room struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Winner     string    `json:"winner"`
	CreatedAt  time.Time `json:"created_at"`
	FinishedAt time.Time `json:"finished_at"`
}

func CreateUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.PostForm("name")
		publicKey := c.PostForm("public_key")

		_, err := db.Exec("INSERT INTO users (name, public_key) VALUES (?, ?)", name, publicKey)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

}

func ListUser(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rows, err := db.Query("SELECT id, name, public_key, created_at FROM users")
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var users []User
		for rows.Next() {
			var user User
			if err := rows.Scan(&user.ID, &user.Name, &user.PublicKey, &user.CreatedAt); err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				return
			}
			users = append(users, user)
		}

		ctx.JSON(200, users)
	}
}
