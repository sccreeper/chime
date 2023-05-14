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
