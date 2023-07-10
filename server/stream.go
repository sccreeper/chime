package main

// Code for handling uploading, downloading & streaming of tracks

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type stream_param struct {
	TrackID string `uri:"track_id" binding:"required"`
}

func handle_stream(ctx *gin.Context) {

	// Verify user
	verified, request_body := verify_user(ctx.Request)
	if !verified {
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

type download_track_query struct {
	ID string `uri:"track_id" binding:"required"`
}

// Download track route, this is also used for caching/downloading for online listening.
func handle_download_track_original(ctx *gin.Context) {

	// Verify user
	verified, request_body := verify_user(ctx.Request)
	if !verified {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	var query download_track_query

	if err := ctx.ShouldBindUri(&query); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var track track_model

	user_id, err := strconv.ParseInt(request_body.UserID, 16, 64)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	track_id, err := strconv.ParseInt(query.ID, 16, 64)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	database.Table(table_tracks).Select("owner", "original").Where("id = ?", track_id).First(&track)

	if user_id != track.Owner {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	ctx.FileAttachment(fmt.Sprintf("/var/lib/chime/tracks/%s/%s", query.ID, track.Original), track.Original)

}

// TOOD: Download album or playlist
func handle_download_playlist_original(ctx *gin.Context) {

}
