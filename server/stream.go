package main

// Code for handling uploading, downloading & streaming of tracks

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type stream_param struct {
	Format string `uri:"format" binding:"required"`
	Track  string `uri:"track_id" binding:"required"`
}

type download_param struct {
	Track string `uri:"track_id" binding:"required"`
}

func handle_stream(ctx *gin.Context) {

	var stream stream_param

	if err := ctx.ShouldBindUri(&stream); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
	}

	// TODO: Stream file

}

// TODO: Download individual track
func handle_download_track(ctx *gin.Context) {

}

// TOOD: Download album or playlist
func handle_download_playlist(ctx *gin.Context) {

}

// TODO: Download individual track
func handle_download_track_original(ctx *gin.Context) {

}

// TOOD: Download album or playlist
func handle_download_playlist_original(ctx *gin.Context) {

}
