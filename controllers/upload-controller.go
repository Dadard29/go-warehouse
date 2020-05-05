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
	queryParam    = "q"

	titleParam = "title"
	artistParam = "artist"
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
func FileGetListLastAdded(w http.ResponseWriter, r *http.Request) {
	accessToken := auth.ParseApiKey(r, accessTokenKey, true)
	if !checkToken(accessToken, w) {
		return
	}

	l, err := managers.FileDbListLastManager()
	if err != nil {
		logger.Error(err.Error())
		api.Api.BuildErrorResponse(http.StatusInternalServerError, "error listing musics", w)
		return
	}

	api.Api.BuildJsonResponse(true, "files listed", l, w)
}

// GET
// Authorization: 	token
// Params: 			None
// Body: 			None

// get the list of albums
func FileGetListAlbums(w http.ResponseWriter, r *http.Request) {
	accessToken := auth.ParseApiKey(r, accessTokenKey, true)
	if !checkToken(accessToken, w) {
		return
	}

	l, err := managers.FileDbListAlbumManager()
	if err != nil {
		logger.Error(err.Error())
		api.Api.BuildErrorResponse(http.StatusInternalServerError, "failed to get album list", w)
		return
	}

	api.Api.BuildJsonResponse(true, "album list retrieved", l, w)

}

// GET
// Authorization: 	token
// Params: 			None
// Body: 			None

// get the list of artist
func FileGetListArtists(w http.ResponseWriter, r *http.Request) {
	accessToken := auth.ParseApiKey(r, accessTokenKey, true)
	if !checkToken(accessToken, w) {
		return
	}

	l, err := managers.FileDbListArtistManager()
	if err != nil {
		logger.Error(err.Error())
		api.Api.BuildErrorResponse(http.StatusInternalServerError, "failed to get artist list", w)
		return
	}

	api.Api.BuildJsonResponse(true, "artist list retrieved", l, w)

}

// DELETE
// Authorization: 	token
// Params: 			title, album, artist
// Body: 			None

// remove file from DB and FS
// todo: remove directories when they get empty after deletion
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

// GET
// Authorization: 	token
// Params: 			None
// Body: 			None

// get a file object from DB
func FileGet(w http.ResponseWriter, r *http.Request) {
	accessToken := auth.ParseApiKey(r, accessTokenKey, true)
	if !checkToken(accessToken, w) {
		return
	}

	title := r.URL.Query().Get(titleParam)
	artist := r.URL.Query().Get(artistParam)

	if title == "" || artist == "" {
		api.Api.BuildMissingParameter(w)
		return
	}

	m, err := managers.FileDbGet(title, artist)
	if err != nil {
		logger.Error(err.Error())
		api.Api.BuildErrorResponse(http.StatusNotFound, "failed to get the music", w)
		return
	}

	api.Api.BuildJsonResponse(true, "music retrieved", m, w)
}

// GET
// Authorization: 	JWT
// Params: 			queryParam
// Body: 			None
func FileSearch(w http.ResponseWriter, r *http.Request) {
	accessToken := auth.ParseApiKey(r, accessTokenKey, true)
	if !checkToken(accessToken, w) {
		return
	}

	q := r.URL.Query().Get(queryParam)
	if q == "" {
		api.Api.BuildMissingParameter(w)
		return
	}

	l, err := managers.FileDbSearchManager(q)
	if err != nil {
		logger.Error(err.Error())
		api.Api.BuildErrorResponse(http.StatusInternalServerError, "error search for musics", w)
		return
	}

	api.Api.BuildJsonResponse(true, "search performed", l, w)

}
