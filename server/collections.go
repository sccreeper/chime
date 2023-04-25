package main

// Handles methods for creating collections and getting collection info.

import (
	"encoding/json"
	"fmt"
	"math"
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
	Radios    []single_collection `json:"radios"`
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
	verified, request_body := verify_user(ctx.Request)
	if !verified {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	var collections = []playlist_model{}
	var radios = []radio_model{}

	owner_id, err := strconv.ParseInt(request_body.UserID, 16, 64)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	database.Table(table_playlists).
		Select("name", "is_album", "id").
		Where("owner = ?", owner_id).Find(&collections)
	database.Table(table_radio).Select("name", "id").Where("owner = ?", owner_id).Find(&radios)

	response := collection_list_response{}

	response.Albums = make([]single_collection, 0)
	response.Playlists = make([]single_collection, 0)
	response.Radios = make([]single_collection, 0)

	for _, v := range collections {

		if v.IsAlbum == 1 {
			response.Albums = append(response.Albums, single_collection{Name: v.Name, ID: strconv.FormatInt(v.ID, 16)})
		} else {
			response.Playlists = append(response.Playlists, single_collection{Name: v.Name, ID: strconv.FormatInt(v.ID, 16)})
		}

	}

	for _, v := range radios {

		response.Radios = append(response.Radios, single_collection{Name: v.Name, ID: strconv.FormatInt(v.ID, 16)})

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
	verified, _ := verify_user(ctx.Request)
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
	verified, request_body := verify_user(ctx.Request)
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
	Size         int    `json:"size"`
}

func handle_get_track_metadata(ctx *gin.Context) {

	// Verify user

	verified, _ := verify_user(ctx.Request)
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
		Size:         int(track.Size),
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

type get_radio_query struct {
	ID string `uri:"radio_id" binding:"required"`
}

type get_radio_response struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Cover       string `json:"cover_id"`
}

func handle_get_radio(ctx *gin.Context) {

	// Verify user & request

	verified, r := verify_user(ctx.Request)
	if !verified {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	var query get_radio_query

	if err := ctx.ShouldBindUri(&query); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	radio_id, err := strconv.ParseInt(query.ID, 16, 64)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	owner_id, err := strconv.ParseInt(r.UserID, 16, 64)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var radio radio_model
	database.Table(table_radio).Select("*").Where("id = ?", radio_id).First(&radio)

	if radio.Owner != owner_id {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	resp := get_radio_response{
		Name:        radio.Name,
		Description: radio.Description,
		Cover:       strconv.FormatInt(radio.CoverID, 16),
		URL:         radio.URL,
	}

	resp_bytes, err := json.Marshal(resp)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Data(http.StatusOK, gin.MIMEJSON, resp_bytes)

}

type create_playlist_query struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CustomCover bool   `json:"custom_cover"`
	IsAlbum     bool   `json:"is_album"`
}

// Create a single collection
func handle_add_collection(ctx *gin.Context) {

	// Verify user & request

	verified, r := verify_user(ctx.Request)
	if !verified {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	if err := ctx.Request.ParseMultipartForm(int64(math.Pow10(6) * 5)); err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Unable to parse form"))
		return
	}

	var playlist create_playlist_query

	if err := json.Unmarshal([]byte(ctx.PostForm("data")), &playlist); err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Invalid request body"))
		return
	}

	user_id, err := strconv.ParseInt(r.UserID, 16, 64)
	if err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: User ID invalid"))
		return
	}

	id := generate_id(table_playlists)

	var is_album int64

	if playlist.IsAlbum {
		is_album = 1
	} else {
		is_album = 0
	}

	// Handle cover

	var cover_id int64 = 0

	if playlist.CustomCover {

		cover_id := generate_id(table_covers)

		cover_file, err := ctx.FormFile("cover")
		if err != nil {
			ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: No cover in form"))
			return
		}

		if err := ctx.SaveUploadedFile(cover_file, fmt.Sprintf("/var/lib/chime/covers/%s", strconv.FormatInt(cover_id, 16))); err != nil {
			ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("500: Unable to save cover"))
			return
		}

		database.Table(table_covers).Create(&cover_model{
			ID:      cover_id,
			AlbumID: 0,
			Owner:   user_id,
		})

	}

	// Create collection record finally.

	database.Table(table_playlists).Create(&playlist_model{
		ID:          id,
		Owner:       user_id,
		Name:        playlist.Name,
		Description: playlist.Description,
		Cover:       cover_id,
		IsAlbum:     is_album,
	})

}

type add_to_collection_query struct {
	TrackID      string `json:"track_id"`
	CollectionID string `json:"collection_id"`
}

// Add tracks to a collection
func add_to_collection(ctx *gin.Context) {

	// Verify user & request

	verified, r := verify_user(ctx.Request)
	if !verified {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	user_id, err := strconv.ParseInt(r.UserID, 16, 64)
	if err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Invalid user ID"))
		return
	}

	var query add_to_collection_query

	if err := ctx.ShouldBindJSON(&query); err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Form formatted incorrectly")) //TODO: Move these to constants. Also only print long message if debug. TLDR; Shorten this boilerplate in some way.
		return
	}

	track_id, err := strconv.ParseInt(query.TrackID, 16, 64)
	if err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Invalid track ID"))
		return
	}

	collection_id, err := strconv.ParseInt(query.CollectionID, 16, 64)
	if err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Invalid collection ID"))
		return
	}

	var collection playlist_model
	var count int64

	// See if track & playlist exists
	database.Table(table_tracks).Where("id = ? AND owner = ?", track_id, user_id).Count(&count)
	database.Table(table_playlists).Where("id = ? AND owner = ?", collection_id, user_id).Count(&count)

	if count == 0 {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Invalid ID "))
		return
	}

	database.Table(table_playlists).Select("*").Where("id = ? AND owner = ?", collection_id, user_id).First(&collection)

	tracks := strings.Split(collection.Tracks, ",")
	tracks = append(tracks, query.TrackID)

	database.Table(table_playlists).Model(&collection).Updates(&playlist_model{Tracks: strings.Join(tracks, ",")})

	ctx.Data(http.StatusOK, gin.MIMEPlain, []byte{})

}
