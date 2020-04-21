package controllers

import (
	"github.com/Dadard29/go-api-utils/auth"
	"github.com/Dadard29/go-warehouse/api"
	"github.com/Dadard29/go-warehouse/managers"
	"github.com/Dadard29/go-warehouse/models"
	"net/http"
)

func DownloadGet(w http.ResponseWriter, r *http.Request) {
	accessToken := auth.ParseApiKey(r, accessTokenKey, true)
	if !checkToken(accessToken, w) {
		return
	}

	title := r.URL.Query().Get("title")
	artist := r.URL.Query().Get("artist")
	album := r.URL.Query().Get("album")

	if title == "" || artist == "" || album == "" {
		api.Api.BuildErrorResponse(http.StatusBadRequest, "missing parameter", w)
		return
	}

	p, err := managers.DownloadGetManager(models.Tags{
		Title:  title,
		Artist: artist,
		Album:  album,
	})

	if err != nil {
		logger.Error(err.Error())
		api.Api.BuildErrorResponse(http.StatusInternalServerError, "error getting file", w)
		return
	}

	http.ServeFile(w, r, p)

}
