package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

func handle_change_password(ctx *gin.Context) {

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
