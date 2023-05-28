package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type search_response struct {
	Tracks []struct {
		ID       string  `json:"id"`
		AlbumID  string  `json:"album_id"`
		Artist   string  `json:"artist"`
		Title    string  `json:"title"`
		Duration float64 `json:"duration"`
		Cover    string  `json:"cover"`
	} `json:"tracks"`

	Collections []struct {
		ID      string `json:"id"`
		Title   string `json:"title"`
		Cover   string `json:"cover"`
		IsAlbum bool   `json:"is_album"`
	} `json:"collections"`

	Radios []struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Cover string `json:"cover"`
	} `json:"radios"`
}

type search_query struct {
	Query string `json:"query"`
}

func handle_search(ctx *gin.Context) {

	// Verify user

	verified, r := verify_user(ctx.Request)
	if !verified {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	// Verify request

	user_id, err := strconv.ParseInt(r.UserID, 16, 64)
	if err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Invalid user ID"))
		return
	}

	var query search_query
	if err := ctx.BindJSON(&query); err != nil {
		ctx.Data(http.StatusBadRequest, gin.MIMEPlain, []byte("400: Invalid request"))
		return
	}

	// Conduct request

	var tracks []track_model
	var collections []playlist_model
	var radios []radio_model

	database.Table(table_tracks).Select("*").Where("owner = ? AND (artist LIKE ? OR name LIKE ?)", user_id, "%"+query.Query+"%", "%"+query.Query+"%").Find(&tracks)
	database.Table(table_playlists).Select("*").Where("owner = ? AND (name LIKE ?)", user_id, "%"+query.Query+"%").Find(&collections)
	database.Table(table_radio).Select("*").Where("owner = ? AND (name LIKE ?)", user_id, "%"+query.Query+"%").Find(&radios)

	resp := search_response{}

	for _, v := range tracks {
		resp.Tracks = append(resp.Tracks, struct {
			ID       string  `json:"id"`
			AlbumID  string  `json:"album_id"`
			Artist   string  `json:"artist"`
			Title    string  `json:"title"`
			Duration float64 `json:"duration"`
			Cover    string  `json:"cover"`
		}{
			ID:       strconv.FormatInt(v.ID, 16),
			AlbumID:  strconv.FormatInt(v.AlbumID, 16),
			Artist:   v.Artist,
			Title:    v.Name,
			Duration: v.Duration,
			Cover:    strconv.FormatInt(v.Cover, 16),
		})
	}

	for _, v := range collections {

		var album bool

		if v.IsAlbum == 1 {
			album = true
		} else {
			album = false
		}

		resp.Collections = append(resp.Collections, struct {
			ID      string `json:"id"`
			Title   string `json:"title"`
			Cover   string `json:"cover"`
			IsAlbum bool   `json:"is_album"`
		}{
			ID:      strconv.FormatInt(v.ID, 16),
			Title:   v.Name,
			Cover:   strconv.FormatInt(v.Cover, 16),
			IsAlbum: album,
		})

	}

	for _, v := range radios {

		resp.Radios = append(resp.Radios, struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Cover string `json:"cover"`
		}{
			ID:    strconv.FormatInt(v.ID, 16),
			Name:  v.Name,
			Cover: strconv.FormatInt(v.CoverID, 16),
		})

	}

	resp_bytes, err := json.Marshal(resp)
	if err != nil {
		ctx.Data(http.StatusInternalServerError, gin.MIMEPlain, []byte("500: Unable to marshal JSON response"))
		return
	}

	ctx.Data(http.StatusOK, gin.MIMEJSON, resp_bytes)
}
