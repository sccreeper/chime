package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/gzip"
	gzipgin "github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var homepage []byte

// Contains code for API paths and delgates them to seperate methods.

func main() {

	// Load homepage

	homepage_bytes, err := os.ReadFile("/dist/index.html")
	if err != nil {
		panic(err)
	}

	homepage = homepage_bytes

	// Load database

	log.Println("Loading database & config...")

	if _, err := os.Stat("/var/lib/chime/data.db"); errors.Is(err, os.ErrNotExist) {
		var running_in_docker bool = false

		if _, err := os.Stat("/.dockerenv"); err == nil {
			running_in_docker = true
		}

		admin_pass, _ := random_string(pass_chars, 10)

		log.Println("No database or default configuration. Creating default database and configuration.")
		first_time_setup(admin_pass, running_in_docker)
		log.Printf("Default user is 'admin' with password '%s' - It is highly recommended you change these values.\n", admin_pass)
		log.Println("Created default configuration and database. Program will exit and you may want to look over configuration files before starting again.")

		if running_in_docker {
			log.Println("Running in Docker, not exiting after setup...")
		} else {
			os.Exit(0)
		}
	}

	database, _ = gorm.Open(sqlite.Open("/var/lib/chime/data.db"), &gorm.Config{})

	r := gin.Default()

	r.Use(gzipgin.Gzip(gzipgin.DefaultCompression, gzip.WithExcludedPathsRegexs([]string{
		`(^\/api\/stream\/).`,
		`(^\/api\/download\/).`,
		`(\/api/\admin\/download_backup\/).`})))

	//Base method for does server exist or not
	r.GET("/api/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Static("/assets/", "/dist/assets/")

	r.GET("/", func(ctx *gin.Context) {
		ctx.Data(http.StatusOK, "text/html; charset=utf-8", homepage)
	})

	// Streaming

	r.GET("/api/stream/:track_id", handle_stream)

	//Uploading tracks & albums

	r.POST("/api/upload", handle_upload)
	r.POST("/api/add_radio", handle_add_radio)

	// Download methods for streaming locally
	r.GET("/api/download/:track_id")
	r.GET("/api/download_collection/:track_id")

	//Download methods for original files
	r.GET("/api/download_original/:track_id", handle_download_track_original)
	r.GET("/api/download_playlist_original/:track_id")

	//Get track & playlist info
	r.GET("/api/get_collection/:collection_id", handle_get_collection)
	r.GET("/api/get_track_metadata/:track_id", handle_get_track_metadata)
	r.GET("/api/get_collections", handle_get_collections)
	r.GET("/api/get_radio/:radio_id", handle_get_radio)
	r.GET("/api/collection/get_cover/:cover_id", handle_get_cover)

	// Modify playlists & albums (collections)

	r.POST("/api/collection/remove_track")
	r.POST("/api/collection/delete", handle_delete_collection)
	r.POST("/api/collection/add", handle_add_collection)
	r.POST("/api/collection/add_collection", handle_add_collection_to_collection)
	r.POST("/api/collection/add_track", add_to_collection)
	r.POST("/api/library/get_track_ids", get_track_ids)

	// Edit endpoints
	r.POST("/api/edit/collection", handle_edit_collection)
	r.POST("/api/edit/reorder_collection", handle_reorder_collection)
	r.POST("/api/edit/track", handle_edit_track)
	r.POST("/api/edit/radio", handle_edit_radio)
	r.POST("/api/edit/favourite", handle_favourite)

	//Search

	r.POST("/api/search", handle_search)

	// Personal stuff

	r.POST("/api/admin/change_username", handle_change_username) //These two are also used for admin user changes.
	r.POST("/api/admin/change_password", handle_change_password)
	r.POST("/api/admin/change_profile_picture")
	r.GET("/api/user/get_profile_picture/:picture_id")
	r.GET("/api/user/get_favourites", handle_get_favourites)

	r.GET("/api/user")

	// Server info & admin

	r.GET("/api/admin/users", handle_get_users)
	r.GET("/api/admin/storage", handle_get_storage)
	r.GET("/api/admin/start_backup", handle_start_backup)
	r.GET("/api/admin/backup_status/:id", handle_get_backup_status)
	r.GET("/api/admin/download_backup/:id", handle_download_backup)
	r.GET("/api/admin/clear_backups", handle_clear_backups)

	r.POST("/api/admin/add_user", handle_add_user)
	r.POST("/api/admin/reset_password", handle_reset_password)
	r.POST("/api/admin/delete_user", handle_delete_user)
	r.POST("/api/admin/toggle_admin", handle_toggle_admin)

	// Embedded
	r.GET("/api/embedded/:session_id/res/:resource_type/:resource_id", handle_embedded_resource)

	// Auth

	r.POST("/api/auth", handle_auth)
	r.GET("/api/auth/session_exists/:session_id", handle_session_exists)

	// Chromecast proxy methods

	r.POST("/api/cast/toggle") //Enables/disables the cast proxy (disabled by default).
	r.GET("/api/cast/enabled")
	r.POST("/api/cast/control") //Play/pause/stop/forward/backwards
	r.POST("/api/cast/set_volume")
	r.GET("/api/cast/get_devices")
	r.GET("/api/cast/get_status/:id")
	r.POST("/api/cast/play_media")

	r.Run("0.0.0.0:80")

}
