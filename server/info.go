package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

// Handles source code for editing playlist, album, track & radio info.

// Favourite tracks

type favourite_request struct {
	ID string `json:"track_id"`
}

type get_favourites_resp struct {
	Favorites []string `json:"favourites"`
}

func handle_favourite(ctx *gin.Context) {

	// Verify user & extract request body
	verified, session := verify_user(ctx.Request)
	if !verified {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	var query favourite_request
	err := ctx.BindJSON(&query)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// IDs
	track_id, err := strconv.ParseInt(query.ID, 16, 64)
	if err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Bad ID"))
		return
	}

	user_id, _ := strconv.ParseInt(session.UserID, 16, 64)

	// Check if track exists

	var count int64
	database.Table(table_tracks).Where("id = ? AND owner = ?", track_id, user_id).Count(&count)

	if count == 0 {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: ID does not exist"))
		return
	}

	//Finally pull user info from DB

	var user user_model
	database.Table(table_users).Select("*").Where("id = ?", user_id).First(&user)

	if !slices.Contains(strings.Split(user.Favourites, ","), query.ID) {

		log.Println("adding to favourites")

		favourite_list := append(strings.Split(user.Favourites, ","), query.ID)

		database.Table(table_users).Model(&user).Updates(&user_model{Favourites: strings.Join(favourite_list, ",")})

		data, err := json.Marshal(get_favourites_resp{Favorites: favourite_list})
		if err != nil {
			ctx.Data(http.StatusInternalServerError, gin.MIMEPlain, []byte("500: Unable to marhsal JSON response"))
			return
		}

		ctx.Data(http.StatusOK, gin.MIMEJSON, data)
		return

	} else {
		log.Println("removing from favourites")

		favourite_list := strings.Split(user.Favourites, ",")

		favourite_list = remove_from_array(favourite_list, slices.Index(favourite_list, query.ID))

		database.Table(table_users).Model(&user).Updates(&user_model{Favourites: strings.Join(favourite_list, ",")})

		data, err := json.Marshal(get_favourites_resp{Favorites: favourite_list})
		if err != nil {
			ctx.Data(http.StatusInternalServerError, gin.MIMEPlain, []byte("500: Unable to marhsal JSON response"))
			return
		}

		ctx.Data(http.StatusOK, gin.MIMEJSON, data)
		return

	}

}

func handle_get_favourites(ctx *gin.Context) {

	// Verify user & extract request body
	verified, session := verify_user(ctx.Request)
	if !verified {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	user_id, _ := strconv.ParseInt(session.UserID, 16, 64)

	var user user_model
	database.Table(table_users).Select("favourites").Where("id = ?", user_id).First(&user)

	resp := get_favourites_resp{Favorites: strings.Split(user.Favourites, ",")}

	data, err := json.Marshal(resp)
	if err != nil {
		ctx.Data(http.StatusInternalServerError, gin.MIMEPlain, []byte("500: Unable to marhsal JSON response"))
		return
	}

	ctx.Data(http.StatusOK, gin.MIMEJSON, data)

}

type reorder_collection_query struct {
	CollectionID string   `json:"collection_id"`
	Tracks       []string `json:"tracks"`
}

// Reorder tracks in collection
func handle_reorder_collection(ctx *gin.Context) {

	// Verify user & extract request body
	verified, session := verify_user(ctx.Request)
	if !verified {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	user_id, _ := strconv.ParseInt(session.UserID, 16, 64)

	var query reorder_collection_query

	if err := ctx.ShouldBindJSON(&query); err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Bad request"))
		return
	}

	collection_id, err := strconv.ParseInt(query.CollectionID, 16, 64)
	if err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Invalid collection ID"))
		return
	}

	// See if collection exists
	if !record_exists(table_playlists, collection_id) {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Collection does not exist"))
		return
	}

	var collection playlist_model
	database.Table(table_playlists).Select("*").Where("id = ?", collection_id).First(&collection)

	if collection.Owner != user_id {
		ctx.Data(http.StatusForbidden, gin.MIMEPlain, []byte("403: User does not own collection"))
		return
	}

	// Verify all tracks exists
	for _, v := range query.Tracks {

		if !record_exists(table_tracks, v) {
			ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Track does not exist"))
			return
		}

	}

	database.Table(table_playlists).Model(&collection).Updates(&playlist_model{Tracks: strings.Join(query.Tracks, ",")})

	ctx.Data(http.StatusOK, gin.MIMEPlain, []byte{})

}

type edit_collection_query struct {
	CollectionID string `json:"collection_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	IsAlbum      bool   `json:"is_album"`
}

func handle_edit_collection(ctx *gin.Context) {

	// Verify user & extract request body
	verified, session := verify_user(ctx.Request)
	if !verified {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	var query edit_collection_query
	if err := ctx.ShouldBindJSON(&query); err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Bad request"))
		return
	}

	if !record_exists(table_playlists, query.CollectionID) {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Collection does not exist"))
		return
	}

	// Verify ownersip of record

	user_id, _ := strconv.ParseInt(session.UserID, 16, 64)
	collection_id, _ := strconv.ParseInt(query.CollectionID, 16, 64)

	var collection playlist_model
	database.Table(table_playlists).Select("*").Where("id = ?", collection_id).First(&collection)

	if collection.Owner != user_id {
		ctx.Data(http.StatusForbidden, gin.MIMEPlain, []byte("403: User does not own record"))
		return
	}

	if collection.Protected == 1 {
		ctx.Data(http.StatusForbidden, gin.MIMEPlain, []byte("403: Record is protected"))
		return
	}

	// Finally apply changes

	var is_album int8
	if query.IsAlbum {
		is_album = 1
	} else {
		is_album = 0
	}

	database.Table(table_playlists).Model(&collection).Updates(&playlist_model{
		Name:        query.Name,
		Description: query.Description,
		IsAlbum:     is_album,
	})

	ctx.Data(http.StatusOK, gin.MIMEPlain, []byte{})

}

type edit_track_query struct {
	TrackID  string `json:"track_id"`
	Name     string `json:"name"`
	Released int64  `json:"released"`
	Artist   string `json:"artist"`
	AlbumID  string `json:"album_id"`
}

func handle_edit_track(ctx *gin.Context) {
	// Verify user & extract request body
	verified, session := verify_user(ctx.Request)
	if !verified {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	var query edit_track_query
	if err := ctx.ShouldBindJSON(&query); err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Bad request"))
		return
	}

	if !record_exists(table_tracks, query.TrackID) {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Track does not exist"))
		return
	}

	// Verify ownership of record

	user_id, _ := strconv.ParseInt(session.UserID, 16, 64)
	collection_id, _ := strconv.ParseInt(query.TrackID, 16, 64)

	var track track_model
	database.Table(table_tracks).Select("*").Where("id = ?", collection_id).First(&track)

	if track.Owner != user_id {
		ctx.Data(http.StatusForbidden, gin.MIMEPlain, []byte("403: User does not own record"))
		return
	}

	// Finally apply changes

	if !record_exists(table_playlists, query.AlbumID) {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Collection does not exist"))
		return
	}

	album_id, _ := strconv.ParseInt(query.AlbumID, 16, 64)

	database.Table(table_tracks).Model(&track).Updates(&track_model{
		Name:     query.Name,
		Released: query.Released,
		Artist:   query.Artist,
		AlbumID:  album_id,
	})

	ctx.Data(http.StatusOK, gin.MIMEPlain, []byte{})

}

type edit_radio_query struct {
	RadioID string `json:"radio_id"`
	Name    string `json:"name"`
	URL     string `json:"url"`
}

func handle_edit_radio(ctx *gin.Context) {
	// Verify user & extract request body
	verified, session := verify_user(ctx.Request)
	if !verified {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	var query edit_radio_query
	if err := ctx.ShouldBindJSON(&query); err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Bad request"))
		return
	}

	if !record_exists(table_radio, query.RadioID) {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Radio does not exist"))
		return
	}

	user_id, _ := strconv.ParseInt(session.UserID, 16, 64)
	radio_id, _ := strconv.ParseInt(query.RadioID, 16, 64)

	var radio radio_model
	database.Table(table_radio).Select("*").Where("id = ?", radio_id).First(&radio)

	if radio.Owner != user_id {
		ctx.Data(http.StatusForbidden, gin.MIMEPlain, []byte("403: User does not own record"))
		return
	}

	database.Table(table_radio).Model(&radio).Updates(&radio_model{
		Name: query.Name,
		URL:  query.URL,
	})

	ctx.Data(http.StatusOK, gin.MIMEPlain, []byte{})

}
