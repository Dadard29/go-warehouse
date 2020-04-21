package repositories

import (
	"errors"
	"fmt"
	"github.com/Dadard29/go-warehouse/models"
	"github.com/bogem/id3v2"
	"github.com/h2non/filetype"
	"io/ioutil"
	"os"
	"path"
)

// create placeholder if needed
func checkPlaceholder(token string) error {
	path2check := path.Join(baseDirStore, token)
	f, err := os.Stat(path2check)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path2check, 0755); err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	}

	if !f.IsDir() {
		return errors.New(fmt.Sprintf("placeholder %s is a file", token))
	}

	return nil
}

func checkArtist(artist string) error {
	path2check := path.Join(baseDirStore, artist)
	f, err := os.Stat(path2check)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path2check, 0755); err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	}

	if !f.IsDir() {
		return errors.New(fmt.Sprintf("artist placeholder %s is a file", artist))
	}

	return nil
}

func checkAlbum(artist string, album string) error {
	path2check := path.Join(baseDirStore, artist, album)
	f, err := os.Stat(path2check)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path2check, 0755); err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	}

	if !f.IsDir() {
		return errors.New(fmt.Sprintf("album placeholder %s is a file", album))
	}

	return nil
}

func CheckFileAudio(path string) bool {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return false
	}

	return filetype.IsAudio(buf)
}

func ReadTags(path string) (models.Tags, error) {
	var fallback models.Tags
	f, err := id3v2.Open(path, id3v2.Options{Parse: true})
	if err != nil {
		return fallback, err
	}

	if f.Title() == "" {
		return fallback, errors.New("title tag empty")
	}

	if f.Artist() == "" {
		return fallback, errors.New("artist tag empty")
	}

	if f.Album() == "" {
		return fallback, errors.New("album tag empty")
	}

	if f.Genre() == "" {
		return fallback, errors.New("genre tag empty")
	}

	if f.Year() == "" {
		return fallback, errors.New("year tag empty")
	}

	return models.Tags{
		Title:       f.Title(),
		Artist:      f.Artist(),
		Album:       f.Album(),
		PublishedAt: f.Year(),
		Genre:       f.Genre(),
	}, nil
}

// return true if file exist
func checkFileExist(tags models.Tags) bool {
	file2check := getFullFilePath(tags)
	_, err := os.Stat(file2check)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		} else if os.IsExist(err) {
			return true
		}
	}

	return true
}
