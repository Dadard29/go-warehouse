package controllers

import (
	"github.com/Dadard29/go-api-utils/auth"
	"github.com/Dadard29/go-warehouse/api"
	"github.com/Dadard29/go-warehouse/managers"
	"github.com/Dadard29/go-warehouse/models"
	"net/http"
)

// GET
// Authorization: 	token
// Params: 			None
// Body: 			None
func FileGetList(w http.ResponseWriter, r *http.Request) {
	accessToken := auth.ParseApiKey(r, accessTokenKey, true)
	if !checkToken(accessToken, w) {
		return
	}

	l, err := managers.FileListManager(accessToken)
	if err != nil {
		logger.Error(err.Error())
		api.Api.BuildErrorResponse(http.StatusInternalServerError, "error listing files", w)
		return
	}

	api.Api.BuildJsonResponse(true, "files listed", l, w)
}

// DELETE
// Authorization: 	token
// Params: 			title, album, artist
// Body: 			None
func FileDelete(w http.ResponseWriter, r *http.Request) {
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

	f, err := managers.FileDeleteManager(accessToken, models.Tags{
		Title:  title,
		Artist: artist,
		Album:  album,
	})

	if err != nil {
		logger.Error(err.Error())
		api.Api.BuildErrorResponse(
			http.StatusInternalServerError, "error while deleting file", w)
		return
	}

	api.Api.BuildJsonResponse(true, "file deleted", f, w)
}

// POST
// Authorization: 	token
// Params: 			None
// Body: 			file to upload
func FileUpload(w http.ResponseWriter, r *http.Request) {
	accessToken := auth.ParseApiKey(r, accessTokenKey, true)
	if !checkToken(accessToken, w) {
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		logger.Error(err.Error())
		api.Api.BuildErrorResponse(
			http.StatusBadRequest, "error parsing form", w)
		return
	}

	file, headers, err := r.FormFile("file")
	if err != nil {
		logger.Error(err.Error())
		api.Api.BuildErrorResponse(
			http.StatusBadRequest, "error getting file", w)
		return
	}

	fileDb, err := managers.FileStoreManager(accessToken, file, headers)
	if err != nil {
		logger.Error(err.Error())
		api.Api.BuildErrorResponse(
			http.StatusInternalServerError, "error storing file", w)
		return
	}

	api.Api.BuildJsonResponse(
		true, "file stored", fileDb, w)
}
