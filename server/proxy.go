package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func handle_proxy(ctx *gin.Context, path_prefix string) {

	remote, err := url.Parse("http://web:3000")
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)

	ctx.Request.Header.Set("Origin", "http://chime:8042")

	proxy.Director = func(r *http.Request) {
		r.Header = ctx.Request.Header
		r.Host = remote.Host
		r.URL.Host = remote.Host
		r.URL.Scheme = remote.Scheme
		r.URL.Path = path_prefix + ctx.Param("proxy")
	}

	proxy.ServeHTTP(ctx.Writer, ctx.Request)

}
