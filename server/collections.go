package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type collection_query struct {
	Session session `json:"session" binding:"required"`
}

type single_collection struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type collection_response struct {
	Albums    []single_collection `json:"albums"`
	Playlists []single_collection `json:"playlists"`
}

func handle_get_collections(ctx *gin.Context) {

	var request_body single_upload

	ctx.BindJSON(&request_body)

	if !verify_user(request_body.Session.ID, request_body.Session.UserID) {
		ctx.AbortWithStatus(http.StatusForbidden)
	}

	var collections = []playlist_model{}

	owner_id, err := strconv.ParseInt(request_body.Session.UserID, 16, 64)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	database.Table("playlists").
		Select("name", "is_album", "id").
		Where("owner = ?", owner_id).Find(&collections)

	fmt.Print("")

	response := collection_response{}

	response.Albums = make([]single_collection, 0)
	response.Playlists = make([]single_collection, 0)

	for _, v := range collections {

		if v.IsAlbum == 1 {

			response.Albums = append(response.Albums, single_collection{Name: v.Name, ID: v.ID})

		} else {
			response.Playlists = append(response.Playlists, single_collection{Name: v.Name, ID: v.ID})
		}

	}

	resp_json, err := json.Marshal(response)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	} else {
		ctx.Data(http.StatusOK, gin.MIMEJSON, resp_json)
	}

}
