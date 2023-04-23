package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
func verify_user(request http.Request) (bool, session) {

	if session_cookie, err := request.Cookie("session"); err != nil {
		fmt.Println(err.Error())
		return false, session{}
	} else {

		// "Escapes" base64

		escaped, _ := base64.StdEncoding.DecodeString(
			strings.ReplaceAll(
				strings.ReplaceAll(
					strings.ReplaceAll(session_cookie.Value, ".", "="),
					"_", "+"),
				"/", "-"))

		var s session
		json.Unmarshal(escaped, &s)

		if _, ok := sessions[s.ID]; !ok {
			return false, session{}
		} else if sessions[s.ID].UserID != s.UserID {
			return false, session{}
		} else {
			return true, s
		}
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
		return
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
				return
			}

			ctx.Data(http.StatusOK, gin.MIMEJSON, response_bytes)
			return

		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"status": "incorrect",
			})
			return
		}

	}

}
