package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/scrypt"
)

type session struct {
	ID     string `json:"session_id" binding:"required"`
	Time   int64  `json:"time" binding:"required"`
	UserID string `json:"user_id" binding:"required"`
}

type login_repsonse struct {
	Status string `json:"status"`
	ID     string `json:"session_id" binding:"required"`
	Time   int64  `json:"time" binding:"required"`
	UserID string `json:"user_id" binding:"required"`
}

var sessions map[string]session

func init() {
	sessions = make(map[string]session)
}

// See if user ID matches session and if session actually exists.
func verify_user(session_id string, user_id string) bool {

	fmt.Println(session_id)
	fmt.Println(user_id)

	var keys []string

	for k, _ := range sessions {
		keys = append(keys, k)
	}

	fmt.Println(keys)

	if _, ok := sessions[session_id]; !ok {
		fmt.Println("doesn't exist")

		return false
	} else if sessions[session_id].UserID != user_id {
		fmt.Println("incorrect user id")

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
				UserID: strconv.FormatInt(user.ID, 16),
			}

			var response login_repsonse = login_repsonse{
				Status: "correct",
				ID:     session_id,
				Time:   sessions[session_id].Time,
				UserID: strconv.FormatInt(user.ID, 16),
			}

			response_bytes, err := json.Marshal(response)
			if err != nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)
			}

			ctx.Data(http.StatusOK, gin.MIMEJSON, response_bytes)

		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"status": "incorrect",
			})
		}

	}

}
