package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Contains code for API paths and delgates them to seperate methods.

func main() {

	// Load database

	log.Println("Loading database & config...")

	if _, err := os.Stat("/var/lib/chime/data.db"); errors.Is(err, os.ErrNotExist) {
		var running_in_docker bool = false

		if _, err := os.Stat("/.dockerenv"); err == nil {
			running_in_docker = true
		}

		admin_pass := generate_password(10)

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

	//Base method for does server exist or not
	r.GET("/api/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Streaming

	r.GET("/api/stream/:format/:track_id")

	//Uploading tracks & albums

	r.POST("/api/upload", handle_upload)

	// Download methods for streaming locally
	r.GET("/api/download/:track_id")
	r.GET("/api/download_playlist/:track_id")

	//Download methods for original files
	r.GET("/api/download_original/:track_id")
	r.GET("/api/download_playlist_original/:track_id")

	//Get track & playlist info
	r.GET("/api/album")
	r.GET("/api/track")

	// Modify playlists & albums

	r.POST("/api/playlist/remove")
	r.POST("/api/playlist/add")
	r.POST("/api/playlist/move")
	r.POST("/api/playlist/change_title")
	r.POST("/api/playlist/change_cover")

	//Search

	r.POST("/api/search")

	// Personal stuff

	r.POST("/api/user/change_username")
	r.POST("/api/user/change_password")
	r.POST("/api/user/change_profile_picture")
	r.POST("/api/user/favourite")

	r.GET("/api/user")

	// Server info & admin

	r.GET("/api/server/users")
	r.GET("/api/server/storage")
	r.GET("/api/server/backup")

	r.POST("/api/admin/new_user")
	r.POST("/api/admin/remove_user")
	r.POST("/api/admin/make_admin")

	// Auth

	r.POST("/api/auth", handle_auth)

	//TODO: Evaluate how Chromecast works
	// Chromecast methods

	r.Run("0.0.0.0:80")

}
