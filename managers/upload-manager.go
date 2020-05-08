package managers

import (
	"errors"
	"fmt"
	"github.com/Dadard29/go-warehouse/models"
	"github.com/Dadard29/go-warehouse/repositories"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
)

const (
	mimeMp3 = "audio/mpeg"

	maxMegaBytes = 10
	maxSize      = maxMegaBytes << (10 * 2)

	maxFilesNumber = 10

	searchLimit = 10
)

func cleanTempFile(path string) {
	err := os.Remove(path)
	if err != nil {
		logger.Error(err.Error())
	}
}

// check if fs and db is same
func FileFsCheck() (bool, error) {
	fsList, err := FileListManager()
	if err != nil {
		return false, err
	}

	dbList := repositories.MusicList()

	if len(fsList) != len(dbList) {
		return false, errors.New("conflicts found")
	}

	for _, f := range fsList {
		check := false
		for _, d := range dbList {
			if f.Metadata.Title == d.Title && f.Metadata.Artist == d.Artist {
				check = true
				break
			}
		}

		if !check {
			return false, errors.New("conflicts found")
		}
	}

	return true, nil
}

// fs
func FileListManager() ([]models.File, error) {
	var f []models.File

	files, err := repositories.ListFiles()
	if err != nil {
		logger.Error("error listing files")
		return f, err
	}

	return files, err
}

func FileDeleteManager(tags models.Tags) (models.File, error) {
	var f models.File

	fileDeleted, err := repositories.RemoveFile(tags)
	if err != nil {
		logger.Error(err.Error())
		return f, errors.New("error while deleting file")
	}

	return fileDeleted, nil
}

func FileStoreManager(file multipart.File, headers *multipart.FileHeader, mp models.MusicParam) (models.File, error) {
	var f models.File

	defer file.Close()

	// check size
	if headers.Size > maxSize {
		return f, errors.New(fmt.Sprintf(
			"file too big: maximum allowed is %d Mb", maxMegaBytes))
	}

	// check mime
	if headers.Header.Get("Content-Type") != mimeMp3 {
		logger.Info(headers.Header.Get("Content-Type"))
		return f, errors.New("bad mime")
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		logger.Error("error reading file")
		return f, err
	}

	tempFileName := "tempfile.mp3"
	tempFilePath := path.Join(repositories.Tmp, tempFileName)
	err = ioutil.WriteFile(tempFilePath, data, 0644)
	if err != nil {
		logger.Error("error writing file")
		return f, err
	}

	// check is audio
	if !repositories.CheckFileAudio(tempFilePath) {
		cleanTempFile(tempFilePath)

		msg := "not an audio file"
		logger.Error(msg)
		return f, errors.New(msg)
	}

	// check if mp3 by reading ID3V2 tag
	tags, err := repositories.ReadTags(tempFilePath)
	if err != nil {
		cleanTempFile(tempFilePath)

		logger.Error(err.Error())
		return f, errors.New("error reading id3v2 tags")
	}


	var fileAdded models.File
	if fileAdded, err = repositories.AddFile(tempFilePath, tags); err != nil {
		cleanTempFile(tempFilePath)

		logger.Error(err.Error())
		return f, errors.New("error storing file in library")
	}

	return fileAdded, nil
}

// db
func FileDbCreateManager(token string, m models.MusicParam, t models.Tags) (models.MusicDto, error) {
	var f models.MusicDto

	mEntity, err := repositories.MusicCreate(token, m, t)
	if err != nil {
		return f, err
	}

	return mEntity.ToDto(), nil
}

func FileDbDelete(title string, artist string) (models.MusicDto, error) {
	var f models.MusicDto

	m, err := repositories.MusicDelete(title, artist)
	if err != nil {
		return f, err
	}

	return m.ToDto(), nil
}

func FileDbListLastManager() ([]models.MusicDto, error) {
	lEntities, err := repositories.MusicListLimit()
	if err != nil {
		return nil, err
	}

	var lDtos = make([]models.MusicDto, 0)
	for _, v := range lEntities {
		lDtos = append(lDtos, v.ToDto())
	}

	return lDtos, nil
}

func FileDbListAlbumManager() ([]models.AlbumDto, error) {
	songList := repositories.MusicList()
	albumList := repositories.MusicAlbumsList()

	var res = make([]models.AlbumDto, 0)
	var imageUrl string
	var artist string

	for _, a := range albumList {
		var titleList = make([]string, 0)
		for _, s := range songList {
			if s.Album == a.Album {
				titleList = append(titleList, s.Title)
				imageUrl = s.ImageUrl
				artist = s.Artist
			}
		}

		res = append(res, models.AlbumDto{
			Name:      a.Album,
			TitleList: titleList,
			Artist:    artist,
			ImageURL:  imageUrl,
		})
	}

	return res, nil
}

func FileDbListArtistManager() ([]models.ArtistDto, error) {
	albumList, err := FileDbListAlbumManager()
	if err != nil {
		return nil, err
	}

	artistList := repositories.MusicArtistsList()


	var res = make([]models.ArtistDto, 0)
	for _, ar := range artistList {
		var arAlbumList = make([]models.AlbumDto, 0)
		for _, al := range albumList {
			if al.Artist == ar.Artist {
				arAlbumList = append(arAlbumList, al)
			}
		}

		res = append(res, models.ArtistDto{
			Name:      ar.Artist,
			AlbumList: arAlbumList,
		})
	}

	return res, nil
}

func FileDbSearchManager(q string) ([]models.MusicDto, error) {
	var lDtos = make([]models.MusicDto, 0)

	for _, field := range []string{repositories.SearchFieldTitle, repositories.SearchFieldArtist, repositories.SearchFieldAlbum} {

		l, err := repositories.MusicSearch(q, field)
		if err != nil {
			return nil, err
		}

		for _, v := range l {
			vDto := v.ToDto()

			// check if v already in lDtos
			check := false
			for _, e := range lDtos {
				if e == vDto {
					check = true
				}
			}

			if !check {
				lDtos = append(lDtos, vDto)
				if len(lDtos) >= searchLimit {
					return lDtos, nil
				}
			}
		}
	}

	return lDtos, nil

}

func FileDbGet(title string, artist string) (models.MusicDto, error) {
	var f models.MusicDto

	m, err := repositories.MusicGetFromTitle(title, artist)
	if err != nil {
		return f, err
	}

	return m.ToDto(), nil
}
