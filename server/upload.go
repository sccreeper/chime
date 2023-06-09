package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dhowden/tag"
	"github.com/gin-gonic/gin"
	"gopkg.in/vansante/go-ffprobe.v2"
)

func handle_upload(ctx *gin.Context) {

	// Parse form and verify user

	var request_body session

	ctx.Request.ParseMultipartForm(int64(math.Pow10(8) * 2.5))

	// Verify user
	verified, request_body := verify_user(ctx.Request)
	if !verified {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
	// Save file & upload IDs.

	file, _ := ctx.FormFile("file")
	filename := path.Base(file.Filename)

	track_id := generate_id(table_tracks)

	id_hex := strconv.FormatInt(track_id, 16)

	os.Mkdir(fmt.Sprintf("/var/lib/chime/tracks/%s/", id_hex), os.ModeDir)
	ctx.SaveUploadedFile(file, fmt.Sprintf("/var/lib/chime/tracks/%s/%s", id_hex, filename))

	// Execute database operations

	f, err := os.Open(fmt.Sprintf("/var/lib/chime/tracks/%s/%s", id_hex, filename))
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	f_info, err := f.Stat()
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	f_size := f_info.Size()

	metadata, err := tag.ReadFrom(f)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	probe_ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	probe_data, err := ffprobe.ProbeURL(probe_ctx, fmt.Sprintf("/var/lib/chime/tracks/%s/%s", id_hex, filename))
	if err != nil {
		ctx.Data(http.StatusInternalServerError, gin.MIMEPlain, []byte("500: Failed to probe duration!"))
	}

	owner_id, err := strconv.ParseInt(request_body.UserID, 16, 64)
	if err != nil {
		log.Println("Failed to user id string to int!")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// If album doesn't exist then add to default "Unsorted" collection.

	if metadata.Title() == "" {
		database.Table(table_tracks).Create(&track_model{
			ID:       track_id,
			Name:     filename,
			Artist:   "Unknown",
			AlbumID:  1,
			Owner:    owner_id,
			Original: filename,
			Size:     f_size,
			Duration: probe_data.Format.DurationSeconds,
			Released: 1968, //Year 2001 a Space Oddesey was released.
		})

		var collection playlist_model
		database.Table(table_playlists).Select("*").Where("id = ?", 1).First(&collection)

		tracks := strings.Split(collection.Tracks, ",")
		tracks = append(tracks, strconv.FormatInt(track_id, 16))

		database.Table(table_playlists).Model(&collection).Updates(&playlist_model{Tracks: strings.Join(tracks, ",")})

	} else {
		// Assume, possibly incorrectly, that the rest of the metadata is present.
		// Create album, or upload to existing if possible.

		var album_title string = metadata.Album()
		var count int64

		user_id, err := strconv.ParseInt(request_body.UserID, 16, 64)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}

		database.Table(table_playlists).Where("name = ? AND is_album = 1 AND owner = ?", album_title, user_id).Count(&count)

		if count == 0 {

			var album_id int64 = generate_id(table_playlists)
			var cover_id int64 = generate_id(table_covers)
			hex_id := strconv.FormatInt(track_id, 16)

			// Create cover record

			if metadata.Picture() == nil {
				cover_id = 0
			} else {
				os.WriteFile(fmt.Sprintf("/var/lib/chime/covers/%s", strconv.FormatInt(cover_id, 16)), metadata.Picture().Data, 0666)

				database.Table(table_covers).Create(&cover_model{
					ID:      cover_id,
					AlbumID: album_id,
					Owner:   owner_id,
				})
			}

			database.Table(table_playlists).Create(&playlist_model{
				ID:      album_id,
				Name:    metadata.Album(),
				IsAlbum: 1,
				Cover:   cover_id,
				Tracks:  hex_id,
				Dates:   strconv.Itoa(metadata.Year()),
				Owner:   owner_id,
			})

			database.Table(table_tracks).Create(&track_model{
				ID:       track_id,
				Name:     metadata.Title(),
				Artist:   metadata.Artist(),
				AlbumID:  album_id,
				Cover:    cover_id,
				Owner:    owner_id,
				Original: filename,
				Size:     f_size,
				Duration: probe_data.Format.DurationSeconds,
				Released: int64(metadata.Year()),
			})

		} else {

			var collection playlist_model
			var new_track_list string

			database.Table(table_playlists).Where("name = ? AND is_album = 1", album_title).First(&collection)

			if len(strings.Split(collection.Tracks, ",")) == 0 {
				new_track_list += strconv.FormatInt(track_id, 16)
			} else {
				new_track_list += fmt.Sprintf("%s,%s", collection.Tracks, strconv.FormatInt(track_id, 16))
			}

			track_position, _ := metadata.Track()

			// Create track record

			database.Table(table_tracks).Create(&track_model{
				ID:       track_id,
				Name:     metadata.Title(),
				Artist:   metadata.Artist(),
				AlbumID:  collection.ID,
				Cover:    collection.Cover,
				Owner:    owner_id,
				Original: filename,
				Size:     f_size,
				Duration: probe_data.Format.DurationSeconds,
				Released: int64(metadata.Year()),
				Position: int64(track_position),
			})

			// Update album list & reorder for album

			var track_list []track_model

			database.Table(table_tracks).Select("*").Where("album_id = ?", collection.ID).Find(&track_list)

			sort.Sort(by_position(track_list))

			var new_track_string string

			new_track_string += strconv.FormatInt(track_list[0].ID, 16)
			for i := 1; i < len(track_list); i++ {

				new_track_string += "," + strconv.FormatInt(track_list[i].ID, 16)

			}

			database.Table(table_playlists).Model(&collection).Updates(&playlist_model{Tracks: new_track_list})

		}

	}

	ctx.Data(http.StatusOK, "text/plain", []byte{})

}

type add_radio_query struct {
	Name string `json:"name" binding:"required"`
	URL  string `json:"url" binding:"required"`
}

func handle_add_radio(ctx *gin.Context) {

	var request_body session

	// Verify user
	verified, request_body := verify_user(ctx.Request)
	if !verified {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	// Parse request

	var query add_radio_query

	if err := ctx.ShouldBindJSON(&query); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}

	radio_id := generate_id(table_radio)

	user_id, err := strconv.ParseInt(request_body.UserID, 16, 64)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	database.Table(table_radio).Create(&radio_model{
		ID:    radio_id,
		URL:   query.URL,
		Name:  query.Name,
		Owner: user_id,

		CoverID:     0,
		Description: "",
	})

	ctx.Data(http.StatusOK, "text/plain", []byte{})

}
