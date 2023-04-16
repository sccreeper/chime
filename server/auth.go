package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/scrypt"
)

type session struct {
	ID     string `json:"id" binding:"required"`
	Time   int64  `json:"time" binding:"required"`
	UserID int64  `json:"user_id" binding:"required"`
}

var sessions map[string]session

func init() {
	sessions = make(map[string]session)
}

// See if user ID matches session and if session actually exists.
func verify_user(session_id string, user_id int64) bool {

	if _, ok := sessions[session_id]; !ok {
		return false
	} else if sessions[session_id].UserID != user_id {
		return false
	} else {
		return true
	}

}

// Handles authentication
func handle_auth(ctx *gin.Context) {

	var username string = ctx.PostForm("u")
	var password string = ctx.PostForm("p")
	var count int64

	// See if user exists
	database.Table(table_users).Where("username = ?", username).Count(&count)

	if count == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "incorrect",
		})
	} else {

		var user user_model
		database.Table(table_users).Where("username = ?", username).First(&user)

		// Hash password & compare

		hash, _ := scrypt.Key([]byte(password), user.Salt, 1<<15, 8, 1, 64)

		if base64.StdEncoding.EncodeToString(hash) == user.Password {

			fmt.Println("Hello World")

			// Create new session

			var session_id string
			var found bool

			for !found {

				session_id = uuid.NewString()

				if _, exists := sessions[session_id]; !exists {
					found = true
				}

			}

			sessions[session_id] = session{
				ID:     session_id,
				Time:   time.Now().Unix(),
				UserID: user.ID,
			}

			ctx.JSON(http.StatusOK, gin.H{
				"status":     "correct",
				"session_id": session_id,
				"time":       sessions[session_id].Time,
				"user_id":    user.ID,
			})
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"status": "incorrect",
			})
		}

	}

}
