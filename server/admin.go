package main

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/sys/unix"
)

const username_check string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890_-"
const new_password_length = 8

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
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Invalid request body"))
		return
	}

	// Conduct request

	user_id_change, err := strconv.ParseInt(query.UserID, 16, 64)
	if err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Invalid user ID"))
		return
	}

	// Check if username already exists
	var count int64
	database.Table(table_users).Where("username = ?", query.Username).Count(&count)
	if count != 0 {
		ctx.JSON(http.StatusOK, gin.H{"status": "username_exists"})
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

		// Check if user is admin and has perm to change other usernames.

		var user user_model
		database.Table(table_users).Select("*").Where("id = ?", user_id).Find(&user)

		if user.IsAdmin != 1 {
			ctx.Data(http.StatusForbidden, gin.MIMEPlain, []byte("403: Insufficient permissions"))
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

type reset_password_query struct {
	UserID string `json:"user_id"`
}

// Resets password to randomized string.
func handle_reset_password(ctx *gin.Context) {

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

	var query reset_password_query
	if err := ctx.ShouldBindJSON(&query); err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Invalid request body"))
		return
	}

	// Check if user is admin and has perm to reset other user passwords.

	var user user_model
	database.Table(table_users).Select("*").Where("id = ?", user_id).Find(&user)

	if user.IsAdmin != 1 {
		ctx.Data(http.StatusForbidden, gin.MIMEPlain, []byte("403: Insufficient permissions"))
		return
	}

	// Generate password & salt

	password, _ := random_string(pass_chars, 10)
	var salt uint64 = random.Uint64()
	salt_bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(salt_bytes, salt)

	hash, err := scrypt.Key([]byte(password), salt_bytes, 1<<15, 8, 1, 64)
	if err != nil {
		ctx.Data(http.StatusInternalServerError, gin.MIMEPlain, []byte("500: Error hashing password"))
		return
	}

	user_change_id, err := strconv.ParseInt(query.UserID, 16, 64)
	if err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Invalid user ID"))
		return
	}

	var user_change user_model
	database.Table(table_users).Select("*").Where("id = ?", user_change_id).First(&user_change)

	database.Table(table_users).Model(&user_change).Updates(&user_model{
		Password: base64.StdEncoding.EncodeToString(hash),
		Salt:     salt_bytes,
	})

	// Finally return the new password

	ctx.JSON(http.StatusOK, gin.H{"password": password})

}

// User admin actions

type get_users_resp struct {
	Users []struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		IsAdmin  bool   `json:"is_admin"`
	} `json:"users"`
}

func handle_get_users(ctx *gin.Context) {

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

	var user user_model
	database.Table(table_users).Select("is_admin").Where("id = ?", user_id).First(&user)

	if user.IsAdmin != 1 {
		ctx.Data(http.StatusForbidden, gin.MIMEPlain, []byte("403: User is not admin"))
		return
	}

	var user_list []user_model
	database.Table(table_users).Select("*").Find(&user_list)

	resp := get_users_resp{
		Users: []struct {
			ID       string `json:"id"`
			Username string `json:"username"`
			IsAdmin  bool   `json:"is_admin"`
		}{},
	}

	for _, v := range user_list {
		var admin bool
		if v.IsAdmin == 1 {
			admin = true
		}

		resp.Users = append(resp.Users, struct {
			ID       string `json:"id"`
			Username string `json:"username"`
			IsAdmin  bool   `json:"is_admin"`
		}{
			ID:       strconv.FormatInt(v.ID, 16),
			Username: v.Username,
			IsAdmin:  admin,
		})
	}

	resp_bytes, err := json.Marshal(resp)
	if err != nil {
		ctx.Data(http.StatusInternalServerError, gin.MIMEPlain, []byte("500: Unable to marshal response"))
		return
	}

	ctx.Data(http.StatusOK, gin.MIMEJSON, resp_bytes)

}

type user_query struct {
	Username string `json:"username"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

func handle_add_user(ctx *gin.Context) {

	// Verify user & request
	verified, session := verify_user(ctx.Request)
	if !verified {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	admin_id, err := strconv.ParseInt(session.UserID, 16, 64)
	if err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Invalid user ID"))
		return
	}

	var admin user_model
	database.Table(table_users).Select("is_admin").Where("id = ?", admin_id).First(&admin)

	if admin.IsAdmin != 1 {
		ctx.Data(http.StatusForbidden, gin.MIMEPlain, []byte("403: User is not admin"))
		return
	}

	// Verify request body
	var query user_query
	if err := ctx.ShouldBindJSON(&query); err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Bad request"))
		return
	}

	// Hash password & create salt, user, etc.
	salt := random.Uint64()
	salt_bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(salt_bytes, salt)
	user_id := generate_id(table_users)

	hash, err := scrypt.Key([]byte(query.Password), salt_bytes, 1<<15, 8, 1, 64)
	if err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("500: Error hashing password"))
		return
	}

	var is_admin int64

	if query.IsAdmin {
		is_admin = 1
	} else {
		is_admin = 0
	}

	// Create database records

	database.Table(table_users).Create(&user_model{
		ID:           user_id,
		Username:     query.Username,
		Password:     base64.StdEncoding.EncodeToString(hash),
		Salt:         salt_bytes,
		IsAdmin:      is_admin,
		SettingsJSON: "",
		Favourites:   "",
	})

	database.Table(table_playlists).Create(&playlist_model{
		ID:          generate_id(table_playlists),
		Name:        "Unsorted",
		IsAlbum:     1,
		Cover:       0,
		Tracks:      "",
		Dates:       "",
		Owner:       user_id,
		Description: "Default album for tracks that don't have any metadata.",
	})

	ctx.Data(http.StatusOK, gin.MIMEPlain, []byte{})

}

type delete_user_query struct {
	UserID string `json:"user_id"`
}

// Deletes a user and all their assets (tracks, playlists, covers etc.)
func handle_delete_user(ctx *gin.Context) {

	// Verify user & request
	verified, session := verify_user(ctx.Request)
	if !verified {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	admin_id, err := strconv.ParseInt(session.UserID, 16, 64)
	if err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Invalid user ID"))
		return
	}

	var admin user_model
	database.Table(table_users).Select("is_admin").Where("id = ?", admin_id).First(&admin)

	if admin.IsAdmin != 1 {
		ctx.Data(http.StatusForbidden, gin.MIMEPlain, []byte("403: User is not admin"))
		return
	}

	var query delete_user_query
	if err := ctx.ShouldBindJSON(&query); err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Invalid query"))
		return
	}

	delete_user_id, err := strconv.ParseInt(query.UserID, 16, 64)
	if err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Invalid query"))
		return
	}

	// Check if user exists
	var count int64
	database.Table(table_users).Where("id = ?", delete_user_id).Count(&count)
	if count == 0 {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Invalid user"))
		return
	}

	// Get user model from database
	var user user_model
	database.Table(table_users).Select("*").Where("id = ?", delete_user_id).First(&user)

	// Find all records and delete them and their associated files.

	var tracks []track_model
	var covers []cover_model
	var collections []playlist_model
	var radios []radio_model

	database.Table(table_tracks).Select("*").Where("owner = ?", delete_user_id).Find(&tracks)
	database.Table(table_covers).Select("*").Where("owner = ?", delete_user_id).Find(&covers)
	database.Table(table_playlists).Select("*").Where("owner = ?", delete_user_id).Find(&collections)
	database.Table(table_radio).Select("*").Where("owner = ?", delete_user_id).Find(&radios)

	for _, v := range tracks {
		os.Remove(fmt.Sprintf("/var/lib/chime/tracks/%s", strconv.FormatInt(v.ID, 16)))
		database.Table(table_tracks).Unscoped().Delete(&v)
	}

	for _, v := range covers {
		os.Remove(fmt.Sprintf("/var/lib/chime/covers/%s", strconv.FormatInt(v.ID, 16)))
		database.Table(table_covers).Unscoped().Delete(&v)
	}

	for _, v := range collections {
		database.Table(table_playlists).Unscoped().Delete(&v)
	}

	for _, v := range radios {
		database.Table(table_radio).Unscoped().Delete(&v)
	}

	// Finally delete user.
	database.Table(table_users).Unscoped().Delete(&user)

	ctx.Data(http.StatusOK, gin.MIMEPlain, []byte{})

}

type toggle_admin_query struct {
	UserID string `json:"user_id"`
}

func handle_toggle_admin(ctx *gin.Context) {

	// Verify user & request
	verified, session := verify_user(ctx.Request)
	if !verified {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	admin_id, err := strconv.ParseInt(session.UserID, 16, 64)
	if err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Invalid user ID"))
		return
	}

	var admin user_model
	database.Table(table_users).Select("is_admin").Where("id = ?", admin_id).First(&admin)

	if admin.IsAdmin != 1 {
		ctx.Data(http.StatusForbidden, gin.MIMEPlain, []byte("403: User is not admin"))
		return
	}

	var query toggle_admin_query
	if err := ctx.ShouldBindJSON(&query); err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Bad request"))
		return
	}

	user_id, err := strconv.ParseInt(query.UserID, 16, 64)
	if err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Bad user ID"))
		return
	}

	// See if user actually exists
	var count int64
	database.Table(table_users).Where("id = ?", user_id).Count(&count)
	if count == 0 {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Invalid user ID"))
		return
	}

	// Check if user id is equal to admin id as admins cannot stop themselves from being admins.
	if user_id == admin_id {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Admins toggle admin on themselves"))
		return
	}

	var user user_model
	database.Table(table_users).Select("*").Where("id = ?", user_id).First(&user)

	if user.IsAdmin == 1 {
		database.Table(table_users).Model(&user).UpdateColumn("is_admin", 0)
	} else {
		database.Table(table_users).Model(&user).UpdateColumn("is_admin", 1)
	}

	ctx.Data(http.StatusOK, gin.MIMEPlain, []byte{})

}

//Other admin actions

type storage_resp struct {
	TotalVolumeSpace int64 `json:"total_volume_space"`
	UsedByOthers     int64 `json:"used_by_others"`
	UsedByChime      int64 `json:"used_by_chime"`
}

func handle_get_storage(ctx *gin.Context) {

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

	var user user_model

	database.Table(table_users).Select("*").Where("id = ?", user_id).First(&user)
	if user.IsAdmin != 1 {
		ctx.Data(http.StatusForbidden, gin.MIMEPlain, []byte("403: User is not admin"))
		return
	}

	// Calculate storage usage
	var size int64

	err = filepath.Walk("/var/lib/chime/", func(path string, info fs.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if !info.IsDir() {
			size += info.Size()
		}

		return err

	})

	if err != nil {
		ctx.Data(http.StatusInternalServerError, gin.MIMEPlain, []byte("500: Internal server error"))
		return
	}

	// Get system storage usage

	var stat unix.Statfs_t
	unix.Statfs("/", &stat)

	var sys_total int64 = int64(stat.Blocks) * stat.Bsize
	var sys_usage int64 = sys_total - (int64(stat.Bavail) * stat.Bsize)

	resp := storage_resp{
		TotalVolumeSpace: sys_total,
		UsedByOthers:     sys_usage - size,
		UsedByChime:      size,
	}

	resp_bytes, err := json.Marshal(resp)
	if err != nil {
		ctx.Data(http.StatusInternalServerError, gin.MIMEPlain, []byte("500: Unable to marshal response"))
		return
	}

	ctx.Data(http.StatusOK, gin.MIMEJSON, resp_bytes)

}

// Backup

type backup struct {
	Progress float64 `json:"progress"`
	ID       int64   `json:"-"`
	Finished bool    `json:"finished"`
	Path     string  `json:"path"`
}

var backups map[int64]backup
var backup_lock = sync.RWMutex{}

func handle_start_backup(ctx *gin.Context) {
	// Verify user & request
	verified, session := verify_user(ctx.Request)
	if !verified {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	admin_id, _ := strconv.ParseInt(session.UserID, 16, 64)

	var admin user_model
	database.Table(table_users).Select("is_admin").Where("id = ?", admin_id).First(&admin)

	if admin.IsAdmin != 1 {
		ctx.Data(http.StatusForbidden, gin.MIMEPlain, []byte("403: User is not admin"))
		return
	}

	// Generate backup id
	var id int64
	for {
		id = random.Int63()
		if _, exists := backups[id]; exists {
			continue
		} else {
			break
		}
	}

	backup_path := fmt.Sprintf("/var/lib/chime/backups/%s", strconv.FormatInt(id, 16))

	file, err := os.Create(backup_path)
	if err != nil {
		ctx.Data(http.StatusInternalServerError, gin.MIMEPlain, []byte("500: Couldn't create backup file"))
	}
	file.Close()

	// Create backup struct

	backups[id] = backup{
		Progress: 0,
		ID:       id,
		Finished: false,
		Path:     backup_path,
	}

	go run_backup(backup_path, id)

	ctx.JSON(http.StatusOK, gin.H{"id": strconv.FormatInt(id, 16)})

}

type backup_status_query struct {
	ID string `uri:"id"`
}

func handle_get_backup_status(ctx *gin.Context) {

	// Verify user & request
	verified, session := verify_user(ctx.Request)
	if !verified {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	admin_id, _ := strconv.ParseInt(session.UserID, 16, 64)

	var admin user_model
	database.Table(table_users).Select("is_admin").Where("id = ?", admin_id).First(&admin)

	if admin.IsAdmin != 1 {
		ctx.Data(http.StatusForbidden, gin.MIMEPlain, []byte("403: User is not admin"))
		return
	}

	var query backup_status_query
	if err := ctx.ShouldBindUri(&query); err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Bad request"))
		return
	}

	backup_id, err := strconv.ParseInt(query.ID, 16, 64)
	if err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Bad ID"))
		return
	}

	backup_lock.RLock()
	defer backup_lock.RUnlock()
	ctx.JSON(http.StatusOK, backups[backup_id])

}

func run_backup(out string, id int64) {

}

func init() {
	backups = make(map[int64]backup, 0)
}
