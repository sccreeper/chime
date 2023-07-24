package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	cast_proxy_root string = "http://host.docker.internal:8080"
)

func handle_cast_control(ctx *gin.Context) {

	// Verify user
	verified, _ := verify_user(ctx.Request)
	if !verified {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	resp, err := http.Post(cast_proxy_root+"/control", gin.MIMEJSON, ctx.Request.Body)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		ctx.AbortWithStatus(resp.StatusCode)
		return
	} else {
		ctx.Status(http.StatusOK)
	}

}

func handle_cast_set_volume(ctx *gin.Context) {
	// Verify user
	verified, _ := verify_user(ctx.Request)
	if !verified {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	resp, err := http.Post(cast_proxy_root+"/set_volume", gin.MIMEJSON, ctx.Request.Body)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		ctx.AbortWithStatus(resp.StatusCode)
		return
	} else {
		ctx.Status(http.StatusOK)
	}
}

func handle_cast_get_devices(ctx *gin.Context) {

	// Verify user
	verified, _ := verify_user(ctx.Request)
	if !verified {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	resp, err := http.Get(cast_proxy_root + "/get_devices")
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	ctx.DataFromReader(http.StatusOK, resp.ContentLength, gin.MIMEJSON, resp.Body, nil)

}

type cast_status_query struct {
	UUID string `uri:"uuid" binding:"required"`
}

func handle_cast_get_status(ctx *gin.Context) {

	// Verify user
	verified, _ := verify_user(ctx.Request)
	if !verified {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	var query cast_status_query

	if err := ctx.ShouldBindUri(&query); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	resp, err := http.Get(fmt.Sprintf("%s/get_status/%s", cast_proxy_root, query.UUID))
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		ctx.AbortWithStatus(resp.StatusCode)
		return
	} else {

		ctx.DataFromReader(http.StatusOK, resp.ContentLength, gin.MIMEJSON, resp.Body, nil)

	}

}

func handle_cast_play_media(ctx *gin.Context) {

	// Verify user
	verified, _ := verify_user(ctx.Request)
	if !verified {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	resp, err := http.Post(cast_proxy_root+"/play_media", gin.MIMEJSON, ctx.Request.Body)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		ctx.AbortWithStatus(resp.StatusCode)
		return
	} else {
		ctx.Status(http.StatusOK)
	}

}
