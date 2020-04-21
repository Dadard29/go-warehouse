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

// check for conflicts between DB and FS
func FileFsCheck(w http.ResponseWriter, r *http.Request) {
	accessToken := auth.ParseApiKey(r, accessTokenKey, true)
	if !checkToken(accessToken, w) {
		return
	}

	c, err := managers.FileFsCheck()
	if err != nil {
		logger.Error(err.Error())
		api.Api.BuildErrorResponse(http.StatusInternalServerError, "conflicts found", w)
		return
	}

	api.Api.BuildJsonResponse(c, "no conflict found", nil, w)

}

// GET
// Authorization: 	token
// Params: 			None
// Body: 			None

// get last added files from DB
func FileGetList(w http.ResponseWriter, r *http.Request) {
	accessToken := auth.ParseApiKey(r, accessTokenKey, true)
	if !checkToken(accessToken, w) {
		return
	}

	l, err := managers.FileDbListManager()
	if err != nil {
		logger.Error(err.Error())
		api.Api.BuildErrorResponse(http.StatusInternalServerError, "error listing musics", w)
		return
	}

	api.Api.BuildJsonResponse(true, "files listed", l, w)
}

// DELETE
// Authorization: 	token
// Params: 			title, album, artist
// Body: 			None

// remove file from DB and FS
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

	_, err := managers.FileDeleteManager(models.Tags{
		Title:  title,
		Artist: artist,
		Album:  album,
	})

	if err != nil {
		logger.Error(err.Error())
		api.Api.BuildErrorResponse(http.StatusInternalServerError, "error while deleting file", w)
		return
	}

	fileDb, err := managers.FileDbDelete(title, artist)
	if err != nil {
		logger.Error(err.Error())
		api.Api.BuildErrorResponse(http.StatusInternalServerError, "error while deleting file in db", w)
		return
	}

	api.Api.BuildJsonResponse(true, "file deleted", fileDb, w)
}

// POST
// Authorization: 	token
// Params: 			None
// Body: 			fileParam, imageUrlParam

// create file in DB and FS
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
		managers.FileDeleteManager(fileStored.Metadata)
		logger.Error(err.Error())
		api.Api.BuildErrorResponse(http.StatusInternalServerError, "error storing file in db", w)
		return
	}

	api.Api.BuildJsonResponse(
		true, "file stored", fileDb, w)
}
