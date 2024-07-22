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

	Session struct {
		ID     string `json:"session_id"`
		Time   int64  `json:"time"`
		UserID string `json:"user_id"`
	} `json:"session"`

	User struct {
		Username string `json:"username"`
		IsAdmin  bool   `json:"is_admin"`
		UserID   string `json:"user_id"`
	} `json:"user"`
}

var sessions map[string]session

func init() {
	sessions = make(map[string]session)
}

// See if user ID matches session and if session actually exists.
func verify_user(request *http.Request) (bool, int64) {

	if session_id, err := request.Cookie("session_id"); err != nil {

		// Legacy for dev purposes atm
		if session_cookie, err := request.Cookie("session"); err != nil {
			fmt.Println(err.Error())
			return false, 0
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
				return false, 0
			} else if sessions[s.ID].UserID != s.UserID {
				return false, 0
			} else {

				user_id, _ := strconv.ParseInt(s.UserID, 16, 64)

				return true, user_id
			}
		}

	} else {

		user_id, err := request.Cookie("user_id")
		if err != nil {
			return false, 0
		}

		if _, ok := sessions[session_id.Value]; !ok {
			return false, 0
		} else if sessions[session_id.Value].UserID != user_id.Value {
			return false, 0
		} else {

			user_id, _ := strconv.ParseInt(user_id.Value, 16, 64)

			return true, user_id
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

		is_admin := false
		if user.IsAdmin == 1 {
			is_admin = true
		}

		if verify_password(&user, password) {

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
				Session: struct {
					ID     string `json:"session_id"`
					Time   int64  `json:"time"`
					UserID string `json:"user_id"`
				}{
					ID:     session_id,
					Time:   sessions[session_id].Time,
					UserID: strconv.FormatInt(user.ID, 16),
				},
				User: struct {
					Username string `json:"username"`
					IsAdmin  bool   `json:"is_admin"`
					UserID   string `json:"user_id"`
				}{
					Username: username,
					IsAdmin:  is_admin,
					UserID:   strconv.FormatInt(user.ID, 16),
				},
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

// Verify whether a session exists, and return the associated user data.

type getUserQuery struct {
	SessionID string `uri:"session_id"`
}

type userResp struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	UserID   string `json:"user_id"`
}

type getUserResponse struct {
	Status string   `json:"status"`
	User   userResp `json:"user"`
}

func getUserFromDb(userId int64) (user user_model, err error) {

	var count int64
	database.Table(table_users).Select("*").Where("id = ?", userId).First(&user).Count(&count)
	if count == 0 {
		err = fmt.Errorf("no user found")
	}
	return

}

func handleGetUser(ctx *gin.Context) {

	var query getUserQuery

	if err := ctx.BindUri(&query); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if _, ok := sessions[query.SessionID]; ok {
		userId, _ := strconv.ParseInt(sessions[query.SessionID].UserID, 16, 64)

		usr, _ := getUserFromDb(userId)
		var isAdmin bool
		if usr.IsAdmin == 1 {
			isAdmin = true
		} else {
			isAdmin = false
		}

		resp, err := json.Marshal(getUserResponse{
			Status: "exists",
			User: userResp{
				UserID:   strconv.FormatInt(userId, 16),
				Username: usr.Username,
				IsAdmin:  isAdmin,
			},
		})

		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.Data(200, gin.MIMEJSON, resp)
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": "doesnt_exist"})
	}

}

// Verifies password matches with user model for other methods that require password authorization
func verify_password(user *user_model, password string) bool {
	hash, _ := scrypt.Key([]byte(password), user.Salt, 1<<15, 8, 1, 64)
	return base64.StdEncoding.EncodeToString(hash) == user.Password
}
