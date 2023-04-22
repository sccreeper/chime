package main

// Code for handling uploading, downloading & streaming of tracks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type stream_param struct {
	TrackID string `uri:"track_id" binding:"required"`
}

func handle_stream(ctx *gin.Context) {

	// Verify user

	var request_body session

	session_json := strings.Join(ctx.Request.Header["Cookie"], "")[len("session="):]
	fmt.Println(session_json)
	json.Unmarshal([]byte(session_json), &request_body)
	if !verify_user(request_body.ID, request_body.UserID) {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	// Get URL params

	var stream stream_param

	if err := ctx.ShouldBindUri(&stream); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	//Check ownership

	var track track_model

	track_id, err := strconv.ParseInt(stream.TrackID, 16, 64)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	owner_id, err := strconv.ParseInt(request_body.UserID, 16, 64)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	database.Table(table_tracks).Select("owner", "original").Where("id = ?", track_id).First(&track)

	if owner_id != track.Owner {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	// After all checks passed, stream track.

	ctx.File(fmt.Sprintf("/var/lib/chime/tracks/%s/%s", stream.TrackID, track.Original))

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
