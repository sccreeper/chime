package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/dhowden/tag"
	"github.com/gin-gonic/gin"
)

type single_upload struct {
	Session session `json:"session" binding:"required"`
}

func handle_upload(ctx *gin.Context) {

	var request_body single_upload

	ctx.BindJSON(&request_body)

	if !verify_user(request_body.Session.ID, request_body.Session.UserID) {
		ctx.AbortWithStatus(http.StatusForbidden)
	}

	file, _ := ctx.FormFile("file")
	id := generate_id("tracks")

	os.Mkdir(fmt.Sprintf("/var/lib/chime/tracks/%d/original", id), os.ModeDir)
	os.Mkdir(fmt.Sprintf("/var/lib/chime/tracks/%d/streamable", id), os.ModeDir)

	ctx.SaveUploadedFile(file, fmt.Sprintf("/var/lib/chime/tracks/%d/original/%s", id, file.Filename))

	// Convert to ogg for streamning on all platforms.

	cmd := exec.Command("ffmpeg", "-i", fmt.Sprintf("/var/lib/chime/tracks/%d/original/%s", id, file.Filename), "-c:a", "copy", fmt.Sprintf("/var/lib/chime/tracks/%d/streamable/%d.ogg", id, id))
	cmd.Run()

	// Execute database operations

	f, err := os.Open(fmt.Sprintf("/var/lib/chime/tracks/%d/original/%s", id, file.Filename))
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	metadata, err := tag.ReadFrom(f)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	owner_id, err := strconv.ParseInt(request_body.Session.UserID, 16, 64)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	if metadata.Title() == "" {
		database.Table(table_tracks).Create(&track_model{
			ID:      generate_id(table_tracks),
			Name:    file.Filename,
			AlbumID: 1,
			Owner:   owner_id,
		})

	} else {
		// Assume, possibly incorrectly, that the rest of the metadata is present.
		// Create album, or upload to existing if possible.

		var album_title string = metadata.Album()
		var count int64

		database.Table(table_playlists).Where("name = ? AND is_album = 1", album_title).Count(&count)

		if count == 0 {

			var album_id int64 = generate_id(table_playlists)
			var track_id int64 = generate_id(table_tracks)

			database.Table(table_playlists).Create(&playlist_model{
				ID:      album_id,
				Name:    metadata.Album(),
				IsAlbum: 1,
				Cover:   0,
				Tracks:  strconv.Itoa(int(track_id)),
				Dates:   strconv.Itoa(metadata.Year()),
				Owner:   owner_id,
			})

			// TODO: Add track to new album.

		} else {

			var track_id int64 = generate_id(table_tracks)
			var playlist playlist_model
			var new_track_list string

			database.Table(table_playlists).Where("name = ? AND is_album = 1", album_title).First(&playlist)

			if len(strings.Split(playlist.Tracks, ",")) == 0 {
				new_track_list += strconv.Itoa(int(track_id))
			} else {
				new_track_list += fmt.Sprintf("%s,%d", playlist.Tracks, track_id)
			}

			database.Table(table_playlists).Model(&playlist).Updates(playlist_model{Tracks: new_track_list})

		}

	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})

}
