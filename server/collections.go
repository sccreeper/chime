package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dhowden/tag"
	"github.com/gin-gonic/gin"
)

type single_collection struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type collection_list_response struct {
	Albums    []single_collection `json:"albums"`
	Playlists []single_collection `json:"playlists"`
}

type collection_query struct {
	AlbumID string `uri:"collection_id" binding"required"`
}

type collection_response struct {
	Title       string           `json:"title"`
	Cover       string           `json:"cover"`
	IsAlbum     bool             `json:"is_album"`
	Tracks      []track_response `json:"tracks"`
	Dates       []string         `json:"dates"`
	Description string           `json:"description"`
}

type track_response struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	AlbumName string `json:"album_name"`
	Released  int    `json:"released"`
	Artist    string `json:"artist"`
	AlbumID   string `json:"album_id"`
	Duration  int    `json:"duration"`
	CoverID   string `json:"cover_id"`
}

func handle_get_collections(ctx *gin.Context) {

	// Verify user
	verified, request_body := verify_user(*ctx.Request)
	if !verified {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	var collections = []playlist_model{}

	owner_id, err := strconv.ParseInt(request_body.UserID, 16, 64)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	database.Table("playlists").
		Select("name", "is_album", "id").
		Where("owner = ?", owner_id).Find(&collections)

	response := collection_list_response{}

	response.Albums = make([]single_collection, 0)
	response.Playlists = make([]single_collection, 0)

	for _, v := range collections {

		if v.IsAlbum == 1 {

			response.Albums = append(response.Albums, single_collection{Name: v.Name, ID: strconv.FormatInt(v.ID, 16)})

		} else {
			response.Playlists = append(response.Playlists, single_collection{Name: v.Name, ID: strconv.FormatInt(v.ID, 16)})
		}

	}

	resp_json, err := json.Marshal(response)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	} else {
		ctx.Data(http.StatusOK, gin.MIMEJSON, resp_json)
	}

}

func handle_get_collection(ctx *gin.Context) {

	// Verify user & extract request body
	verified, _ := verify_user(*ctx.Request)
	if !verified {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	var query_params collection_query
	err := ctx.BindUri(&query_params)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Query database for playlist

	var collection playlist_model

	album_id, err := strconv.ParseInt(query_params.AlbumID, 16, 64)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	database.Table("playlists").Select("*").
		Where("id = ?", album_id).
		First(&collection)

	var response_struct collection_response = collection_response{}

	// Add values
	response_struct.Title = collection.Name
	response_struct.Cover = strconv.FormatInt(collection.Cover, 16)
	response_struct.Description = collection.Description

	if collection.IsAlbum == 1 {
		response_struct.IsAlbum = true
	} else {
		response_struct.IsAlbum = false
	}

	response_struct.Dates = strings.Split(collection.Dates, ",")

	// Add tracks

	response_struct.Tracks = make([]track_response, 0)

	for _, v := range strings.Split(collection.Tracks, ",") {

		var track track_model
		var track_collection playlist_model

		id, err := strconv.ParseInt(v, 16, 64)
		if err != nil {
			continue
		}

		database.Table("tracks").Select("*").Where("id = ?", id).First(&track)
		database.Table("playlists").Select("name").Where("id = ?", track.AlbumID).First(&track_collection)

		response_struct.Tracks = append(response_struct.Tracks, track_response{
			ID:        strconv.FormatInt(id, 16),
			Name:      track.Name,
			Released:  int(track.Released),
			AlbumName: track_collection.Name,
			Artist:    track.Artist,
			AlbumID:   strconv.FormatInt(track.AlbumID, 16),
			Duration:  int(track.Duration),
			CoverID:   strconv.FormatInt(track.AlbumID, 16),
		})

	}

	// Serialize and return bytes.

	resp_bytes, err := json.Marshal(response_struct)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	} else {
		ctx.Data(http.StatusOK, gin.MIMEJSON, resp_bytes)
	}

}

type get_cover_query struct {
	CoverID string `uri:"cover_id" binding:"required"`
}

func handle_get_cover(ctx *gin.Context) {

	// Verify user
	verified, request_body := verify_user(*ctx.Request)
	if !verified {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
	// Get cover to query

	var query get_cover_query
	ctx.ShouldBindUri(&query)

	cover_id, err := strconv.ParseInt(query.CoverID, 16, 64)
	if err != nil {
		fmt.Println(query.CoverID)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	user_id, err := strconv.ParseInt(request_body.UserID, 16, 64)
	if err != nil {
		fmt.Println(request_body.UserID)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var cover cover_model
	database.Table(table_covers).Select("owner").Where("id = ?", cover_id).First(&cover)

	if cover.Owner != user_id {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	} else {
		ctx.File(fmt.Sprintf("/var/lib/chime/covers/%s", query.CoverID))
	}

}

type track_metadata_query struct {
	ID string `uri:"track_id" binding:"required"`
}

type track_metadata_response struct {
	Title        string `json:"title"`
	AlbumName    string `json:"album_name"`
	CoverID      string `json:"cover_id"`
	Artist       string `json:"artist"`
	OriginalFile string `json:"original_file"`
	Format       string `json:"format"`
	Duration     int    `json:"duration"`
	Released     int    `json:"released"`
}

func handle_get_track_metadata(ctx *gin.Context) {

	// Verify user

	// Verify user
	verified, _ := verify_user(*ctx.Request)
	if !verified {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	var query track_metadata_query

	if err := ctx.ShouldBindUri(&query); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	track_id, err := strconv.ParseInt(query.ID, 16, 64)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}

	var track track_model
	database.Table(table_tracks).Select("*").Where("id = ?", track_id).First(&track)

	var response track_metadata_response = track_metadata_response{
		Title:        track.Name,
		CoverID:      strconv.FormatInt(track.Cover, 16),
		Artist:       track.Artist,
		OriginalFile: track.Original,
		Duration:     int(track.Duration),
		Released:     int(track.Released),
	}

	var track_album playlist_model
	database.Table(table_playlists).Select("name").Where("id = ?", track.AlbumID).First(&track_album)
	response.AlbumName = track_album.Name

	track_file, err := os.Open(fmt.Sprintf("/var/lib/chime/tracks/%s/%s", strconv.FormatInt(track.ID, 16), track.Original))
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	metadata, err := tag.ReadFrom(track_file)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	switch metadata.FileType() {
	case tag.FLAC:
		response.Format = "FLAC"
	case tag.MP3:
		response.Format = "MP3"
	case tag.OGG:
		response.Format = "WAV"
	case tag.UnknownFileType:
		response.Format = "Uknown"
	default:
		response.Format = "Unsupported Format"
	}

	response_bytes, err := json.Marshal(response)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Data(http.StatusOK, gin.MIMEJSON, response_bytes)

}
