package main

import (
	"encoding/base64"
	"encoding/binary"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/scrypt"
)

const username_check string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890_-"

type change_username_query struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Allows user to change username or password (non-admin actions)
func handle_change_username(ctx *gin.Context) {

	// Verify user & request
	verified, session := verify_user(ctx.Request)
	if !verified {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	user_id, err := strconv.ParseInt(session.UserID, 16, 64)
	if err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Invalid user ID"))
		return
	}

	var query change_username_query

	if err := ctx.ShouldBindJSON(&query); err != nil {
		ctx.Data(http.StatusForbidden, gin.MIMEPlain, []byte("403: Invalid request body"))
		return
	}

	// Conduct request

	user_id_change, err := strconv.ParseInt(query.UserID, 16, 64)
	if err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Invalid user ID"))
		return
	}

	if user_id_change == user_id {

		var user user_model

		database.Table(table_users).Select("*").Where("id = ?", user_id).Find(&user)

		if !verify_string(query.Username, username_check) {
			ctx.JSON(http.StatusOK, gin.H{"status": "bad_username"})
			return
		} else {

			// Verify password

			if !verify_password(&user, query.Password) {
				ctx.JSON(http.StatusOK, gin.H{"status": "bad_auth"})
				return
			}

			// Finally change username

			database.Table(table_users).Model(&user).Update("username", query.Username)
			ctx.JSON(http.StatusOK, gin.H{"status": "success"})
			return
		}

	} else {

		if !verify_string(query.Username, username_check) {
			ctx.JSON(http.StatusOK, gin.H{"status": "bad_username"})
			return
		}

		// Check if user is admin and has perm to change other user IDs.

		var user user_model
		database.Table(table_users).Select("*").Where("id = ?", user_id).Find(&user)

		if user.IsAdmin != 1 {
			ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("403: Insufficient permissions"))
			return
		}

		// Is admin so continue with change

		var count int64
		var user_change user_model

		database.Table(table_users).Select("*").Where("id = ?", user_id_change).Count(&count)
		if count == 0 {
			ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Invalid user ID"))
			return
		}

		database.Table(table_users).Select("*").Where("id = ?", user_id_change).Find(&user_change)
		database.Table(table_users).Model(&user_change).Update("username", query.Username)

		ctx.Data(http.StatusOK, gin.MIMEPlain, []byte{})
		return

	}

}

type change_password_query struct {
	OldPassword  string `json:"old_password"`
	NewPassword0 string `json:"new_password_0"`
	NewPassword1 string `json:"new_password_1"`
}

// User reset password
func handle_change_password(ctx *gin.Context) {

	// Verify user & request
	verified, session := verify_user(ctx.Request)
	if !verified {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	user_id, err := strconv.ParseInt(session.UserID, 16, 64)
	if err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Invalid user ID"))
		return
	}

	var query change_password_query

	if err := ctx.ShouldBindJSON(&query); err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Invalid request body"))
		return
	}

	// Verify all details of request

	if query.NewPassword0 != query.NewPassword1 {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Non-matching passwords"))
		return
	}

	var user user_model
	var count int64

	if database.Table(table_users).Where("id = ?", user_id).Count(&count); count == 0 {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Invalid user ID"))
		return
	}

	// Finally verify current password and change
	database.Table(table_users).Select("*").Where("id = ?", user_id).Find(&user)

	var salt uint64 = random.Uint64()
	var salt_bytes []byte = make([]byte, 8)
	binary.LittleEndian.PutUint64(salt_bytes, salt)

	hash, _ := scrypt.Key([]byte(query.NewPassword0), salt_bytes, 1<<15, 8, 1, 64)

	database.Table(table_users).Model(&user).Updates(&user_model{Password: base64.StdEncoding.EncodeToString(hash), Salt: salt_bytes})

	ctx.Data(http.StatusOK, gin.MIMEPlain, []byte{})

}

// Admin reset password
func handle_reset_password(ctx *gin.Context) {

}

// User admin actions
func handle_get_users(ctx *gin.Context) {

}

func handle_add_user(ctx *gin.Context) {

}

func handle_remove_user(ctx *gin.Context) {

}

// Other admin actions
func handle_get_backup(ctx *gin.Context) {

}

func handle_get_storage(ctx *gin.Context) {

}
