package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Embedded endpoints
// This exists for "embedded" devices/libaries where things like headers and cookies are not supported.
// Examples include Chromecast devices, operating system media previews/thumbnails, or custom "embedded" environments
// The session ID used is the users session ID returned in the inital auth request.

const (
	cover_resource string = "cover"
	track_resource string = "track"
)

type embedded_token_query struct {
	ID           string `uri:"resource_id"`
	ResourceType string `uri:"resource_type"`
	SessionID    string `uri:"session_id"`
}

func handle_embedded_resource(ctx *gin.Context) {

	// Verify user & request
	var query embedded_token_query
	if err := ctx.ShouldBindUri(&query); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest) //Minimal error codes, simple as possible
		return
	}

	if _, ok := sessions[query.SessionID]; !ok {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user_id, _ := strconv.ParseInt(sessions[query.SessionID].UserID, 16, 64)

	switch query.ResourceType {
	case cover_resource:
		if !record_exists[string](table_covers, query.ID) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		cover_id, _ := strconv.ParseInt(query.ID, 16, 64)

		var cover cover_model
		database.Table(table_covers).Select("owner").Where("id = ?", cover_id).First(&cover)

		if cover.Owner != user_id {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		// Finally return file

		ctx.File(fmt.Sprintf("/var/lib/chime/covers/%s", query.ID))

		return
	case track_resource:
		if !record_exists[string](table_tracks, query.ID) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		track_id, _ := strconv.ParseInt(query.ID, 16, 64)

		var track track_model
		database.Table(table_covers).Select("owner", "original").Where("id = ?", track_id).First(&track)

		if track.Owner != user_id {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		// Finally return file

		ctx.File(fmt.Sprintf("/var/lib/chime/tracks/%s/%s", query.ID, track.Original))

		return

	default:
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

}
