package controllers

import (
	"github.com/Dadard29/go-api-utils/auth"
	"github.com/Dadard29/go-warehouse/api"
	"github.com/Dadard29/go-warehouse/managers"
	"github.com/Dadard29/go-warehouse/models"
	"net/http"
)

const (
	fileParam     = "file"
	imageUrlParam = "image_url"
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

	l, err := managers.FileListManager()
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
// Params: 			(form-data): fileParam, imageUrlParam
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

	imageUrl := r.Form.Get(imageUrlParam)
	m := models.MusicParam{
		ImageUrl: imageUrl,
	}
	if !m.CheckSanity() {
		api.Api.BuildMissingParameter(w)
		return
	}

	file, fileHeaders, err := r.FormFile(fileParam)
	if err != nil {
		logger.Error(err.Error())
		api.Api.BuildErrorResponse(
			http.StatusBadRequest, "error getting file", w)
		return
	}

	// store file
	fileStored, err := managers.FileStoreManager(file, fileHeaders, m)
	if err != nil {
		logger.Error(err.Error())
		api.Api.BuildErrorResponse(
			http.StatusInternalServerError, "error storing file", w)
		return
	}

	// create in db
	fileDb, err := managers.FileDbCreateManager(accessToken, m, fileStored.Metadata)
	if err != nil {
		managers.FileDeleteManager(accessToken, fileStored.Metadata)
		logger.Error(err.Error())
		api.Api.BuildErrorResponse(http.StatusInternalServerError, "error storing file in db", w)
		return
	}

	api.Api.BuildJsonResponse(
		true, "file stored", struct {
			f models.File
			m models.MusicDto
		}{
			fileStored,
			fileDb,
		}, w)
}
